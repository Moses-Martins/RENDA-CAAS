package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"RENDA-CAAS/models"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var GoogleOauthConfig *oauth2.Config

func init() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not loaded")
	}

	// Now build the config
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

}

// Handler to start OAuth2 login
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := GoogleOauthConfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handler to handle the callback from Google
// Handler to handle the callback from Google
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	token, err := GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Look up user by email
	var user models.User
	err = UserCollection.FindOne(ctx, bson.M{"email": userInfo.Email}).Decode(&user)
	if err != nil {
		// User does not exist, create with "User" role for all products
		user = models.User{
			Name:         userInfo.Name,
			Email:        userInfo.Email,
			ProductRoles: map[string]string{"Renda360": "User", "Scale": "User", "Horizon": "User"},
			SuperAdmin:   false,
		}

		result, err := UserCollection.InsertOne(ctx, user)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			user.ID = oid.Hex()
		}
	}

	// Generate JWT (same as your normal login)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"productRoles": user.ProductRoles,
		"exp":          time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"RENDA-CAAS/config"
	"RENDA-CAAS/models"

	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection

func InitUserCollection() {
	UserCollection = config.DB.Collection("users")
}

// registerForProduct is a helper to register a user for a specific product.
func registerForProduct(w http.ResponseWriter, r *http.Request, mainProduct string) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	input.Password = string(hashedPassword)

	// Assign roles: mainProduct gets "User", others get "Viewer"
	allProducts := []string{"renda360", "Scale", "Horizon"}
	input.ProductRoles = make(map[string]string)
	for _, product := range allProducts {
		if product == mainProduct {
			input.ProductRoles[product] = "User"
		} else {
			input.ProductRoles[product] = "Viewer"
		}
	}

	// Check if user exists
	count, _ := UserCollection.CountDocuments(context.TODO(), bson.M{"email": input.Email})
	if count > 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	_, err := UserCollection.InsertOne(context.TODO(), input)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login authenticates a user and returns a JWT token with user info and roles.
func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := UserCollection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		"user_id":      user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"productRoles": user.ProductRoles,
		"exp":          time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func UpdateUserPrivilege(w http.ResponseWriter, r *http.Request) {
	// Get the requester (the one making the change)
	requester, err := getUserFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.PrivilegeUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the target user (the one being updated)
	var targetUser models.User
	err = UserCollection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&targetUser)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	// Superadmin logic
	if requester.SuperAdmin {
		// Superadmin can do anything
		updatePrivilege(&targetUser, req.Product, req.Role)
	} else {
		// Not superadmin: must be admin of the product
		if requester.ProductRoles[req.Product] != "Admin" {
			http.Error(w, "Forbidden: not admin of this product", http.StatusForbidden)
			return
		}
		// Admins cannot promote to admin
		if req.Role == "Admin" {
			http.Error(w, "Forbidden: only superadmin can assign admin role", http.StatusForbidden)
			return
		}
		// Admins cannot demote another admin
		if targetUser.ProductRoles[req.Product] == "Admin" && req.Role != "Admin" {
			http.Error(w, "Forbidden: only superadmin can demote an admin", http.StatusForbidden)
			return
		}
		// Admin can promote/demote to User/Viewer or remove access
		updatePrivilege(&targetUser, req.Product, req.Role)
	}

	// Save changes
	_, err = UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": targetUser.Email},
		bson.M{"$set": bson.M{"productRoles": targetUser.ProductRoles}},
	)
	if err != nil {
		http.Error(w, "Failed to update privilege: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Privilege updated successfully"})
}

// Helper to update the privilege in the user struct
func updatePrivilege(user *models.User, product, role string) {
	if role == "" {
		delete(user.ProductRoles, product)
	} else {
		user.ProductRoles[product] = role
	}
}

// Me returns the current user's info and product access details from the JWT token.
func Me(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return
	}

	// Return user info and product access details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}

func getUserFromToken(r *http.Request) (*models.User, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token")
	}

	var user models.User
	err = UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RegisterRenda360 handles registration for Renda360 product.
func RegisterRenda360(w http.ResponseWriter, r *http.Request) {
	registerForProduct(w, r, "renda360")
}

func RegisterScale(w http.ResponseWriter, r *http.Request) {
	registerForProduct(w, r, "Scale")
}

func RegisterHorizon(w http.ResponseWriter, r *http.Request) {
	registerForProduct(w, r, "Horizon")
}

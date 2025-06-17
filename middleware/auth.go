package middleware

import (
	"RENDA-CAAS/models"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validProducts = map[string]string{
	"renda360": "renda360",
	"scale":    "Scale",
	"horizon":  "Horizon",
}

func AdminOrUserForProduct(userCollection *mongo.Collection) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			routeProduct := strings.ToLower(vars["product"])

			backendProduct, ok := validProducts[routeProduct]
			if !ok {
				http.Error(w, "Invalid product", http.StatusNotFound)
				return
			}

			tokenString := r.Header.Get("Authorization")
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			email, ok := claims["email"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var user models.User
			err = userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			role := user.ProductRoles[backendProduct]
			if role == "Admin" || role == "User" || user.SuperAdmin {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "Forbidden: only Admin or User can access this dashboard", http.StatusForbidden)
		})
	}
}

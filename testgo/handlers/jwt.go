package handlers

import (
	"net/http"
	"strings"
	"user-service/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Define a secret key for JWT
var jwtKey = []byte("your-secret-key")

// CustomClaims represents custom claims for JWT
type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(user models.User) (string, error) {
	claims := CustomClaims{
		Username:       user.Username,
		Password:       user.Password,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthenticateJWT middleware authenticates JWT tokens
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse and verify the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token.Valid {
			// If the token is valid, proceed to the next handler
			context.Set(r, "user", token.Claims)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

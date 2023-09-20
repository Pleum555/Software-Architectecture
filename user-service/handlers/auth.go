package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Pleum555/User-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword generates a salted hash of the given password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// verifyPassword checks if the provided password matches the hashed password
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Register a new user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	fmt.Print(user)
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Generate a salted hash of the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	user.Password = hashedPassword

	_, err1 := userCollection.InsertOne(context.Background(), &user)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error registering user")
		return
	}

	// You can now access the user data from the request body
	fmt.Printf("Received user data: %+v\n", user)

	// Handle user registration logic here
	// Insert the user into the database or perform any other necessary actions

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")
}

// Login and generate a JWT token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var inputUser models.User
	_ = json.NewDecoder(r.Body).Decode(&inputUser)

	// Find the user by username (you should fetch user data from a database)
	var foundUser models.User
	filter := bson.M{"username": inputUser.Username}
	err := userCollection.FindOne(context.Background(), filter).Decode(&foundUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	err1 := VerifyPassword(inputUser.Password, foundUser.Password)
	if err1 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Password does not match the hash:", err)
		return
	}

	tokenString, _ := GenerateJWT(foundUser)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tokenString)
}

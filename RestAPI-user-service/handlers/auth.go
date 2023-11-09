package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Pleum555/User-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// Function to get a user by username from the database
func getUserByUsername(username string) *models.User {
	var user models.User
	filter := bson.M{"username": username}
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // User not found
		}
		fmt.Println("Error querying database:", err)
		return nil // Handle the error as needed
	}
	return &user
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"` // Include the 'token' field here
	Msg     string `json:"msg,omitempty"`   // Include the 'msg' field if needed
	// You can add other fields as needed
}

func sendTokenResponse(user models.User, statusCode int, w http.ResponseWriter) {
	// Create token
	token, err := user.GenerateJWT()
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	options := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 30), // 30 days
		HttpOnly: true,
	}

	if os.Getenv("NODE_ENV") == "production" {
		options.Secure = true
	}

	http.SetCookie(w, options)

	response := RegisterResponse{Success: true, Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	fmt.Print(user)
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Check if the username is already taken
	existingUser := getUserByUsername(user.Username)
	if existingUser != nil {
		w.WriteHeader(http.StatusConflict) // HTTP 409 Conflict
		fmt.Fprintf(w, "Username already exists")
		return
	}

	// Generate a salted hash of the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error hashing password")
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
	w.WriteHeader(http.StatusCreated)
	sendTokenResponse(user, http.StatusOK, w)
}

// Login user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Find the user by username (you should fetch user data from a database)
	var foundUser models.User
	filter := bson.M{"username": request.Username}
	err := userCollection.FindOne(context.Background(), filter).Decode(&foundUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	err1 := VerifyPassword(request.Password, foundUser.Password)
	if err1 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Password does not match the hash:", err)
		return
	}

	tokenString, _ := GenerateJWT(foundUser)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tokenString)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "none",
		Expires:  time.Now().Add(10 * time.Second),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true, "data": {}}`))
}

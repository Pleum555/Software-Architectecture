package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define MongoDB client and collection
var client *mongo.Client
var userCollection *mongo.Collection

// Define a secret key for JWT
var jwtKey = []byte("your-secret-key")

// User struct to represent a user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Users slice to store user data (in-memory storage for this example)
var users = []User{
	{Username: "user1", Password: "password1"},
	{Username: "user2", Password: "password2"},
}

// CustomClaims represents custom claims for JWT
type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// Handlers

// Register a new user
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	json.NewDecoder(r.Body).Decode(&user)

	_, err := userCollection.InsertOne(context.Background(), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error registering user")
		return
	}
	fmt.Println(err)

	// You can now access the user data from the request body
	fmt.Printf("Received user data: %+v\n", user)

	// Handle user registration logic here
	// Insert the user into the database or perform any other necessary actions

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")
}

// Login and generate a JWT token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var inputUser User
	_ = json.NewDecoder(r.Body).Decode(&inputUser)

	// Find the user by username (you should fetch user data from a database)
	var foundUser User
	filter := bson.M{"username": inputUser.Username, "password": inputUser.Password}
	err := userCollection.FindOne(context.Background(), filter).Decode(&foundUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Authentication failed")
		return
	}

	// Create a JWT token
	claims := CustomClaims{
		Username: foundUser.Username,
		Password: foundUser.Password,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating token")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tokenString)
}

// Middleware to authenticate JWT token
func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb+srv://witsaroot:pleum555@suckseat.ipx0cx4.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Initialize the user collection
	userCollection = client.Database("Users").Collection("users")

	r := mux.NewRouter()

	// Register and login routes
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")

	// Protected route (requires authentication)
	r.Handle("/protected", authenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the protected route!")
	})))

	http.Handle("/", r)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

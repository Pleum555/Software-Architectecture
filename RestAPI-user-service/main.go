package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Pleum555/User-service/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("./config.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	handlers.InitMongoDB()
	r := mux.NewRouter()

	// Register and login/logout routes
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")

	// Protect the /getcurrentlocation route with JWT authentication
	r.Handle("/verifyuserdetail", handlers.AuthenticateJWT(http.HandlerFunc(handlers.VerifyUserDetail))).Methods("GET")
	r.Handle("/getuserdetail", handlers.AuthenticateJWT(http.HandlerFunc(handlers.GetUserDetail))).Methods("GET")
	r.Handle("/getcurrentuserlocation", handlers.AuthenticateJWT(http.HandlerFunc(handlers.GetCurrentUserLocation))).Methods("GET")

	// Add the new route and handler for updating user details
	r.Handle("/updateuserdetail", handlers.AuthenticateJWT(http.HandlerFunc(handlers.UpdateUserDetail))).Methods("PUT")

	http.Handle("/", r)

	fmt.Println("Server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Pleum555/User-service/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load("./config.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handlers.InitMongoDB()

	r := mux.NewRouter()

	// Register and login/logout routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Protect the /getcurrentlocation route with JWT authentication
	r.Handle("/verifyuserdetail", handlers.AuthenticateJWT(http.HandlerFunc(handlers.VerifyUserDetail))).Methods("GET")
	r.Handle("/getuserdetail", handlers.AuthenticateJWT(http.HandlerFunc(handlers.GetUserDetail))).Methods("GET")
	r.Handle("/messagetoplaces", handlers.AuthenticateJWT(http.HandlerFunc(handlers.MessageToPlaces))).Methods("GET")

	// Add the new route and handler for updating user details
	r.Handle("/updateuserdetail", handlers.AuthenticateJWT(http.HandlerFunc(handlers.UpdateUserDetail))).Methods("PUT")

	// Create a new CORS middleware handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Change this to the specific origins you want to allow
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true, // Set to true to enable debugging logs
	})

	// Use the CORS middleware with your router
	handler := c.Handler(r)

	http.Handle("/", handler)

	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/Pleum555/User-service/handlers"
	"github.com/gorilla/mux"
)

func main() {
	handlers.InitMongoDB()
	r := mux.NewRouter()

	// Register and login routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Protect the /getcurrentlocation route with JWT authentication
	r.Handle("/getcurrentlocation", handlers.AuthenticateJWT(http.HandlerFunc(handlers.GetCurrentLocation))).Methods("GET")

	http.Handle("/", r)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

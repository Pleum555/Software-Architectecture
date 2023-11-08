package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Pleum555/User-service/models"
	context2 "github.com/gorilla/context"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUserDetail(w http.ResponseWriter, r *http.Request) {
	var updateUser models.User
	_ = json.NewDecoder(r.Body).Decode(&updateUser)

	userValue := context2.Get(r, "user")
	user, _ := userValue.(*CustomClaims)

	checkUser := getUserByUsername(user.Username)

	if updateUser.Username != "" {
		checkUser.Username = updateUser.Username
	}
	if updateUser.Name != "" {
		checkUser.Name = updateUser.Name
	}
	if updateUser.Surname != "" {
		checkUser.Surname = updateUser.Surname
	}
	if updateUser.Tel != "" {
		checkUser.Tel = updateUser.Tel
	}
	if updateUser.Status != "" {
		checkUser.Status = updateUser.Status
	}
	if updateUser.Role != "" {
		checkUser.Role = updateUser.Role
	}
	if updateUser.Password != "" {
		// Generate a salted hash of the password
		hashedPassword, err := HashPassword(updateUser.Password)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error hashing password")
			return
		}
		checkUser.Password = hashedPassword
	}

	// Assuming you have a MongoDB collection named userCollection
	// You can use the `FindOneAndUpdate` method to update the user's details
	// The following code updates the user document with the provided username
	// and sets the new user details from `updateUser`
	filter := bson.M{"username": checkUser.Username}
	update := bson.M{"$set": checkUser}
	err := userCollection.FindOneAndUpdate(context.Background(), filter, update).Decode(&checkUser)
	if err != nil {
		fmt.Println("Error updating database:", err)
		// Handle the error as needed
		// You might want to return an error response to the client here
		return
	}

	userJSON, _ := json.Marshal(checkUser)
	fmt.Fprintf(w, "%s", userJSON)
}

func VerifyUserDetail(w http.ResponseWriter, r *http.Request) {
	userValue := context2.Get(r, "user")
	user, _ := userValue.(*CustomClaims)
	checkUser := getUserByUsername(user.Username)
	// fmt.Println("Username:", user.Username)
	// fmt.Println("Password:", user.Password)

	// userJSON, _ := json.Marshal(checkUser)
	// fmt.Fprintf(w, "%s", userJSON)
	// Create a new struct with only the "role" field
	roleData := struct {
		Role string `json:"role"`
	}{
		Role: string(checkUser.Role),
	}

	roleJSON, err := json.Marshal(roleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(roleJSON)

	return
}

func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	userValue := context2.Get(r, "user")
	user, _ := userValue.(*CustomClaims)
	checkUser := getUserByUsername(user.Username)
	// fmt.Println("Username:", user.Username)
	// fmt.Println("Password:", user.Password)

	userJSON, _ := json.Marshal(checkUser)
	fmt.Fprintf(w, "%s", userJSON)

	return
}

func GetCurrentUserLocation(w http.ResponseWriter, r *http.Request) {
	// Extract the username from the request or wherever you get it
	userValue := context2.Get(r, "user")
	user, _ := userValue.(*CustomClaims)
	// fmt.Println("Username:", user.Username)
	// fmt.Println("Password:", user.Password)
	username := user.Username

	var requestData models.Location
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	location := requestData.Location

	// Create a connection to the RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Replace with your RabbitMQ server connection details
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	queueName := location // Replace with your queue name
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Define the message body (in this case, the username)
	messageBody := []byte(username)

	// Publish the message to the queue
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBody,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	// Send a response to the client
	fmt.Fprintf(w, "Username sent to RabbitMQ queue: %s", location)
	return
}

// ...

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Pleum555/User-service/models"
	"github.com/gorilla/context"
	"github.com/streadway/amqp"
)

var rabbitMQURL = "amqp://guest:guest@localhost:5672/"

func UpdateUserDetail(w http.ResponseWriter, r *http.Request) {
	userValue := context.Get(r, "user")
	user, _ := userValue.(*CustomClaims)
	checkUser := getUserByUsername(user.Username)

	// fmt.Println("Username:", user.Username)
	// fmt.Println("Password:", user.Password)

	userJSON, _ := json.Marshal(checkUser)
	fmt.Fprintf(w, "%s", userJSON)
	return
}

func VerifyUserDetail(w http.ResponseWriter, r *http.Request) {
	userValue := context.Get(r, "user")
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
	userValue := context.Get(r, "user")
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
	userValue := context.Get(r, "user")
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

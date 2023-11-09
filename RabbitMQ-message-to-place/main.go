package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func ReceiveFromQueue(place string) {
	// Connect to RabbitMQ server.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Replace with your RabbitMQ server details.
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel.
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue to consume from.
	_, err = ch.QueueDeclare(
		place, // Queue name
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare the queue: %v", err)
	}

	// Consume messages from the queue.
	msgs, err := ch.Consume(
		place, // Queue name
		"",    // Consumer
		true,  // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume from the queue: %v", err)
	}

	// Process received messages.
	for msg := range msgs {
		message := string(msg.Body)
		fmt.Printf("Received message at place %s: %s\n", place, message)

		// You can process the received message here as needed.
		// For example, update the place based on the received message.
	}
}

func main() {
	fmt.Print("Enter the place: ")
	var place string
	fmt.Scanln(&place)

	ReceiveFromQueue(place)

	// Keep the program running to continue receiving messages.
	select {}
}

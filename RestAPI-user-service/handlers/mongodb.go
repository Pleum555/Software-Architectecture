package handlers

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func InitMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	// fmt.Println(mongoURI)
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to Database")
	}

	userCollection = client.Database("Users").Collection("users")
}

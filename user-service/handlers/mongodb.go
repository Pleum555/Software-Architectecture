package handlers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://witsaroot:pleum555@suckseat.ipx0cx4.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to Database")
	}

	userCollection = client.Database("Users").Collection("users")
}

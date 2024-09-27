package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection string
	connectionString := "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	// Set up a context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Access the pr_analyzer database and repositories collection
	database := client.Database("pr_analyzer")
	collection := database.Collection("repositories")

	// Define the repository document
	repo := bson.M{
		"name": "MetaMask",
		"url":  "https://github.com/MetaMask/metamask-extension.git",
	}

	// Insert the repository into the collection
	result, err := collection.InsertOne(ctx, repo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted repository with ID: %v\n", result.InsertedID)
}

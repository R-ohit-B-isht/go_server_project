package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go_server_project/models"
	"go_server_project/routes"
)

var client *mongo.Client

func main() {
	// Set up file logging
	logFile, err := os.OpenFile("pullrequest_collection.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	log.Println("Server started successfully")
	log.Println("Logging setup verified")
	log.Println("Additional log statement to verify logging functionality")

	// Initialize Bloom filter
	models.InitPRBloomFilter(1000000, 0.01) // Capacity: 1 million, False Positive Rate: 1%
	fmt.Println("Bloom filter initialized")

	// Set up router
	router := mux.NewRouter()

	// Get collections
	db := client.Database("pr_analyzer")
	repoCollection := db.Collection("repositories")
	prCollection := db.Collection("pullrequests")
	clusterCollection := db.Collection("clusters")
	contributorCollection := db.Collection("contributors")
	analysisResultCollection := db.Collection("analysisresults")
	processingJobCollection := db.Collection("processingjobs")
	configurationCollection := db.Collection("configurations")

	// Register routes
	routes.RegisterRepositoryRoutes(router, repoCollection)
	routes.RegisterPullRequestRoutes(router, prCollection, repoCollection)
	routes.RegisterClusterRoutes(router, clusterCollection)
	routes.RegisterContributorRoutes(router, contributorCollection)
	routes.RegisterAnalysisResultRoutes(router, analysisResultCollection)
	routes.RegisterProcessingJobRoutes(router, processingJobCollection)
	routes.RegisterConfigurationRoutes(router, configurationCollection)
	routes.RegisterBulkUploadRoutes(router, prCollection, repoCollection)

	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"},
		AllowCredentials: true,
		Debug: true,
	})

	// Set up HTTP server with CORS middleware
	handler := c.Handler(router)

	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

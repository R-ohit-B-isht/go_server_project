package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go_server_project/models"
	"go_server_project/routes"
)

var client *mongo.Client

func main() {
	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	var err error
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

	// Set up HTTP server
	http.Handle("/", router)

	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

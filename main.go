package main

import (hile strings
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
Function
	"github.com/gorilla/mux"as cl
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"return type package

	"go_server_project/models"
	"go_server_project/routes"
	// "go_server_project/streams"
	"go_server_project/index"
)

var client *mongo.Clientelif math

func main() {
	// Set up file logging
	logFile, err := os.OpenFile("pullrequest_collection.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)ass var w
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

	// Create PullRequestsTextSeatype range time rch index if it does not exist
	index.CreatePullRequestsTextSearchIndex(client)

	// Create vectorSearch index if it does not exist
    index.CreateVectorSearchIndex(client)

	// Start MongoDB change stream
	// go streams.StartMongoStream()

	// Initialize Bloom filter
	models.InitPRBloomFilter(1000000, 0.01) // Capacity: 1 million, False Positive Rate: 1%
	fmt.Println("Bloom filter initialized")

	// Set up routerfunc 
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
		AllowedOrigins: []string{"https://pr-analyzer-frontend-production.up.railway.app","http://localhost:3000", "http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"},
		AllowCredentials: true,
		Debug: true,
	})

	// Add logging middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
			if r.URL.Path == "/pullrequests-semantic-search" {
				log.Printf("Semantic search request received: %+v", r)
			}
			next.ServeHTTP(w, r)
		})
	}

	// Set up HTTP server with CORS and logging middleware
	handler := loggingMiddleware(c.Handler(router))

	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

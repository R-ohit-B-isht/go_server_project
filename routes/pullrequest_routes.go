package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"strings"

	"go_server_project/models"
	"go_server_project/streams"
	"os"
)

func RegisterPullRequestRoutes(router *mux.Router, prCollection *mongo.Collection, repoCollection *mongo.Collection) {
	// Create a unique index on the PRId field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"prId": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := prCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Error creating unique index on PRId: %v", err)
	}

	router.HandleFunc("/pullrequests", createPullRequest(prCollection, repoCollection)).Methods("POST")
	router.HandleFunc("/pullrequests-collect", collectPullRequests(prCollection, repoCollection)).Methods("POST")
	router.HandleFunc("/pullrequests", getPaginatedPullRequests(prCollection)).Methods("GET")
	router.HandleFunc("/pullrequests/{id}", getPullRequest(prCollection)).Methods("GET")
	router.HandleFunc("/pullrequests/{id}", updatePullRequest(prCollection)).Methods("PUT")
	router.HandleFunc("/pullrequests/{id}", deletePullRequest(prCollection)).Methods("DELETE")
	router.HandleFunc("/pullrequests-search", fullTextPullRequestsSearch(prCollection, repoCollection)).Methods("POST")
	router.HandleFunc("/pullrequests-semantic-search", semanticPullRequestsSearch(prCollection, repoCollection)).Methods("POST")
	router.HandleFunc("/pullrequests-sync", syncPullRequests()).Methods("POST")
	router.HandleFunc("/pullrequests-syncLevel", getCurrentRepoSyncLevel(prCollection, repoCollection)).Methods("GET")
}
strings
func semanticPullRequestsSearch(prCollection *mongo.Collection, repoCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Starting semanticPullRequestsSearch function")
		defer log.Println("Ending semanticPullRequestsSearch function")

		var searchRequest struct {
			SearchText string `json:"searchText"`
		}
		if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("Received search request: %+v", searchRequest)

		repoID := r.URL.Query().Get("id")continue
		log.Printf("Repository ID: %s", repoID)
		objectID, err := primitive.ObjectIDFromHex(repoID)
		if err != nil {
			log.Printf("Error converting repository ID: %v", err)
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)parameter
			returnrange
		}
		loDstruct istrig.Printf("Converted repository ID to ObjectID: %s", objectID.Hex())

		// Generate embeddings using OpenAI API
		log.Println("Generating embeddings using OpenAI API")
		client := &http.Client{}
		requestBody := fmt.Sprintf(`{"input": "%s", "model": "text-embedding-ada-002"}`, searchRequest.SearchText)
		log.Printf("OpenAI API Request Body: %s", requestBody)sync range
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", strings.NewReader(requestBody))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")trace
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("OPENAI_API_KEY is not set")
			http.Error(w, "OpenAI API key is not configured", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		log.Println("Sending request to OpenAI API")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request to OpenAI API: %v", err)
			http.Error(w, "Failed to generate embeddings", http.StatusInternalServerError)
	go def os io context continue 		return
		}
		defer resp.Body.Close()
		log.Printf("OpenAI API Response Status: %d", resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			returnbuted else
		}
		log.Printf("OpenAI API Response Body: %s", string(bodyBytes))

		var embeddingResponse struct {
			Data []struct {
				Embedding []float64 `json:"embedding"`
			} `json:"daerrors sync except case ta"`
		}
		if err := json.Unmarshal(bodyBytes, &embeddingResponse); err != nil {
			log.Printf("Error decoding embedding response: %v", err)
			http.Error(w, "Failed to decode embedding response", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully decoded embedding response")
		log.Printf("Embedding response structure: %+v", embeddingResponse)

		if len(embeddingResponse.Data) == 0 {
			log.Println("Error: Embedding response data is empty")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Failed to generate embeddings",
				"documents": []interface{}{},
			})
			return
		}

		log.Printf("Number of embedding data entries: %d", len(embeddingResponse.Data))
		log.Printf("First embedding data entry: %+v", embeddingResponse.Data[0])
		log.Printf("Embedding type: %T, length: %d", embeddingResponse.Data[0].Embedding, len(embeddingResponse.Data[0].Embedding))

		queryVector := embeddingResponse.Data[0].Embedding
		log.Printf("Generated embedding vector with %d dimensions", len(queryVector))

		// Define the pipeline for vector search
		pipeline := mongo.Pipeline{
			{{Key: "$vectorSearch", Value: bson.M{
				"index": "vectorSemanticSearch",
				"path": "embedding",
				"queryVector": queryVector,
				"numCandidates": 100,
				"limit": 10,
			}}},
			{{Key: "$addFields", Value: bson.M{
				"score": bson.M{"$meta": "vectorSearchScore"},
			}}},
			{{Key: "$match", Value: bson.M{
				"repository": objectID,
			}}},
		}
		log.Printf("Defined MongoDB aggregation pipeline: %+v", pipeline)
		log.Printf("Query vector type: %T", queryVector)
		log.Printf("Query vector length: %d", len(queryVector))
		log.Printf("Query vector first 5 elements: %v", queryVector[:5])

		log.Println("Executing MongoDB aggregation pipeline")
		log.Printf("Pipeline details: %+v", pipeline)
		log.Printf("Query vector: %v", queryVector[:5]) // Log first 5 elements of query vector
		cursor, err := prCollection.Aggregate(r.Context(), pipeline)
		if err != nil {
			log.Printf("Error executing search: %v", err)
			log.Printf("Error type: %T", err)
			log.Printf("Full error details: %+v", err)
			http.Error(w, "Failed to execute search", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())
		log.Println("Successfully executed MongoDB aggregation")

		var pullRequests []bson.M
		if err = cursor.All(r.Context(), &pullRequests); err != nil {
			log.Printf("Error decoding search results: %v", err)
			log.Printf("Error type: %T", err)
			log.Printf("Full error details: %+v", err)
			http.Error(w, "Failed to decode search results", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully decoded %d pull requests", len(pullRequests))

		// Remove embedding field from each document
		for i := range pullRequests {
			delete(pullRequests[i], "embedding")
		}

		log.Printf("Found %d pull requests matching the search criteria", len(pullRequests))

		response := struct {
			Success   bool     `json:"success"`
			Count     int      `json:"count"`
			Documents []bson.M `json:"documents"`
		}{
			Success:   true,
			Count:     len(pullRequests),
			Documents: pullRequests,
		}
		log.Printf("Prepared response: %+v", response)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully encoded and sent response")
		log.Println("semanticPullRequestsSearch function completed successfully")
	}
}

func syncPullRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go streams.Synchroniser()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Synchronization started"})
	}
}

func getCurrentRepoSyncLevel(prCollection *mongo.Collection, repoCollection *mongo.Collection) http.HandlerFunc {
	// return the number of PRs for a particular repository with repository as mongo id in the get request field that have embeddings field vs the total number of PRs you'd be given a get request like this  'http://localhost:8080/pullrequests-syncLevel?id=66f70e56d5c8e3c9d8d91252'
	return func(w http.ResponseWriter, r *http.Request) {
		// return a json object with the total number of PRs and the number of PRs with embeddings for now just return 1,2 for total and embeddings respectively


		log.Println("getCurrentRepoSyncLevel function called")
		repoID := r.URL.Query().Get("id")
		log.Printf("Received repoID: %s", repoID)

		objectID, err := primitive.ObjectIDFromHex(repoID)
		if err != nil {
			log.Printf("Error converting repoID to ObjectID: %v", err)
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		filter := bson.M{"repository": objectID}
		log.Printf("Total PRs filter: %+v", filter)
		totalPRs, err := prCollection.CountDocuments(r.Context(), filter)
		if err != nil {
			log.Printf("Error counting total PRs: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Total PRs count: %d", totalPRs)

		embeddingsFilter := bson.M{"repository": objectID, "embedding": bson.M{"$exists": true}}
		log.Printf("Embeddings filter: %+v", embeddingsFilter)
		embeddingsPRs, err := prCollection.CountDocuments(r.Context(), embeddingsFilter)
		if err != nil {
			log.Printf("Error counting PRs with embeddings: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("PRs with embeddings count: %d", embeddingsPRs)

		// Fetch a sample document to verify structure
		sampleDoc := bson.M{}
		err = prCollection.FindOne(r.Context(), filter).Decode(&sampleDoc)
		if err != nil {
			log.Printf("Error fetching sample document: %v", err)
		} else {
			log.Printf("Sample document structure: %+v", sampleDoc)
		}

		response := struct {
			TotalPRs      int64 `json:"totalPRs"`
			EmbeddingsPRs int64 `json:"embeddingsPRs"`
		}{
			TotalPRs:      totalPRs,
			EmbeddingsPRs: embeddingsPRs,
		}

		log.Printf("Total PRs: %d, PRs with embeddings: %d", totalPRs, embeddingsPRs)
		log.Printf("Response struct: %+v", response)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		log.Println("Response sent successfully")
	}
}

func fullTextPullRequestsSearch(prCollection *mongo.Collection, repoCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var searchRequest struct {
			SearchText string `json:"searchText"`
		}
		if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		repoID := r.URL.Query().Get("id")
		objectID, err := primitive.ObjectIDFromHex(repoID)
		if err != nil {
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		pipeline := mongo.Pipeline{
			bson.D{{Key: "$search", Value: bson.M{
				"index": "PullRequestsTextSearch",
				"text": bson.M{
					"query": searchRequest.SearchText,
					"path":  bson.M{"wildcard": "*"},
				},
			}}},
			bson.D{{Key: "$match", Value: bson.M{"repository": objectID}}},
		}
		log.Printf("Search pipeline: %+v", pipeline)
		log.Printf("Search text: %s", searchRequest.SearchText)
		log.Printf("Repository ID: %s", objectID.Hex())

		cursor, err := prCollection.Aggregate(r.Context(), pipeline)
		if err != nil {
			log.Printf("Error executing search: %v", err)
			http.Error(w, "Failed to execute search", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var pullRequests []models.PullRequest
		if err = cursor.All(r.Context(), &pullRequests); err != nil {
			log.Printf("Error decoding search results: %v", err)
			http.Error(w, "Failed to decode search results", http.StatusInternalServerError)
			return
		}
		log.Printf("Search returned %d results", len(pullRequests))

		if len(pullRequests) == 0 {
			log.Printf("No results found")
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}

		response := struct {
			Success   bool                  `json:"success"`
			Count     int                   `json:"count"`
			Documents []models.PullRequest `json:"documents"`
		}{
			Success:   true,
			Count:     len(pullRequests),
			Documents: pullRequests,
		}

		log.Printf("Response: %+v", response)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getPaginatedPullRequests(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		repoID := query.Get("id")
		pageNumber, _ := strconv.Atoi(query.Get("pageNumber"))
		pageSize, _ := strconv.Atoi(query.Get("pageSize"))

		if pageNumber < 1 {
			pageNumber = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}

		skip := (pageNumber - 1) * pageSize

		filter := bson.M{}
		if repoID != "" {
			objectID, err := primitive.ObjectIDFromHex(repoID)
			if err != nil {
				http.Error(w, "Invalid repository ID", http.StatusBadRequest)
				return
			}
			filter["repository"] = objectID
		}

		options := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

		cursor, err := collection.Find(r.Context(), filter, options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var pullRequests []models.PullRequest
		if err = cursor.All(r.Context(), &pullRequests); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(pullRequests) == 0 {
			pullRequests = []models.PullRequest{} // Return empty array instead of null
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pullRequests)
	}
}

func createPullRequest(collection *mongo.Collection, repoCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pullRequest models.PullRequest
		json.NewDecoder(r.Body).Decode(&pullRequest)

		// Check if collection is empty and clear Bloom filter if needed
		count, err := collection.CountDocuments(r.Context(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count == 0 {
			models.ClearPRBloomFilter()
		}

		// Check if PR already exists in Bloom filter
		repo, err := getRepositoryByID(repoCollection, pullRequest.Repository)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if repo.CheckPRBloomFilter(pullRequest.PRId) {
			http.Error(w, "PR may already exist", http.StatusConflict)
			return
		}

		result, err := collection.InsertOne(r.Context(), pullRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update repository's bloom filter
		repo, err = getRepositoryByID(repoCollection, pullRequest.Repository)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add PR to Bloom filter
		repo.AddToPRBloomFilter(pullRequest.PRId)
		// Update repository in database with new bloom filter
		if err := updateRepositoryBloomFilter(repoCollection, repo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the corresponding repository
		prID := result.InsertedID.(primitive.ObjectID)
		updatedRepo, err := getRepositoryByID(repoCollection, pullRequest.Repository)
		if err != nil {
			log.Printf("Error finding repository: %v", err)
		} else {
			err = updatedRepo.DeserializeBloomFilter()
			if err != nil {
				log.Printf("Error deserializing Bloom filter: %v", err)
			}
			_, err = repoCollection.UpdateOne(
				r.Context(),
				bson.M{"_id": pullRequest.Repository},
				bson.M{"$push": bson.M{"pullRequests": prID}},
			)
			if err != nil {
				log.Printf("Error updating repository: %v", err)
			} else {
				err = repo.SerializeBloomFilter()
				if err != nil {
					log.Printf("Error serializing Bloom filter: %v", err)
				}
				_, err = repoCollection.UpdateOne(
					r.Context(),
					bson.M{"_id": pullRequest.Repository},
					bson.M{"$set": bson.M{"serializedBloom": repo.SerializedBloom}},
				)
				if err != nil {
					log.Printf("Error updating serialized Bloom filter: %v", err)
				}
			}
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getRepositoryByID(repoCollection *mongo.Collection, repoID primitive.ObjectID) (*models.Repository, error) {
	var repo models.Repository
	err := repoCollection.FindOne(context.Background(), bson.M{"_id": repoID}).Decode(&repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func updateRepositoryBloomFilter(repoCollection *mongo.Collection, repo *models.Repository) error {
	err := repo.SerializeBloomFilter()
	if err != nil {
		return err
	}
	_, err = repoCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": repo.ID},
		bson.M{"$set": bson.M{"serializedBloom": repo.SerializedBloom}},
	)
	return err
}

func getAllPullRequests(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pullRequests []models.PullRequest

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var pullRequest models.PullRequest
			cursor.Decode(&pullRequest)
			pullRequests = append(pullRequests, pullRequest)
		}

		json.NewEncoder(w).Encode(pullRequests)
	}
}

func getPullRequest(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var pullRequest models.PullRequest
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&pullRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(pullRequest)
	}
}

func updatePullRequest(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var pullRequest models.PullRequest
		json.NewDecoder(r.Body).Decode(&pullRequest)

		update := primitive.M{
			"$set": pullRequest,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(pullRequest)
	}
}

func deletePullRequest(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Pull Request deleted successfully"})
	}
}

func collectPullRequests(prCollection *mongo.Collection, repoCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up file logging
		logFilePath := os.Getenv("PR_LOG_FILE_PATH")
		if logFilePath == "" {
			logFilePath = "pullrequest_collection.log"
		}
		log.Printf("Attempting to create log file at: %s", logFilePath)

		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening log file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer logFile.Close()

		fileLogger := log.New(logFile, "", log.LstdFlags)
		fileLogger.Printf("Log file opened successfully at: %s", logFilePath)
		log.Printf("Log file opened successfully at: %s", logFilePath)

		var requestBody struct {
			StartDate  string `json:"startDate"`
			EndDate    string `json:"endDate"`
			DateFormat string `json:"dateFormat"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			fileLogger.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		fileLogger.Printf("Request body decoded: %+v", requestBody)

		repoID := r.URL.Query().Get("id")
		if repoID == "" {
			fileLogger.Printf("Repository ID is missing")
			http.Error(w, "Repository ID is required", http.StatusBadRequest)
			return
		}
		objectID, err := primitive.ObjectIDFromHex(repoID)
		if err != nil {
			fileLogger.Printf("Invalid repository ID: %s", repoID)
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}
		fileLogger.Printf("Repository ID (hex): %s", objectID.Hex())
		fileLogger.Printf("Repository ID (string): %s", repoID)

		// Find the repository in the repo collection
		var repo models.Repository
		filter := bson.M{"_id": objectID}
		fileLogger.Printf("Querying database with filter: %+v", filter)
		err = repoCollection.FindOne(r.Context(), filter).Decode(&repo)
		if err != nil {
			fileLogger.Printf("Repository not found: %v", err)
			fileLogger.Printf("Query result: %+v", repo)
			// Log all repository IDs in the collection
			cursor, _ := repoCollection.Find(r.Context(), bson.M{})
			var repos []models.Repository
			cursor.All(r.Context(), &repos)
			for _, r := range repos {
				fileLogger.Printf("Existing repository ID: %s", r.ID.Hex())
			}
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}
		fileLogger.Printf("Repository found: %s", repo.URL)

		// Deserialize the Bloom filter
		err = repo.DeserializeBloomFilter()
		if err != nil {
			fileLogger.Printf("Error deserializing Bloom filter: %v", err)
			// Initialize a new Bloom filter if deserialization fails
			repo.InitBloomFilter(1000, 0.01) // Adjust capacity and false positive rate as needed
			err = repo.SerializeBloomFilter()
			if err != nil {
				fileLogger.Printf("Error initializing and serializing new Bloom filter: %v", err)
				http.Error(w, "Error processing repository data", http.StatusInternalServerError)
				return
			}
		}

		// Initialize GitHub client
		ctx := context.Background()
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			fileLogger.Printf("GITHUB_TOKEN environment variable is not set")
			http.Error(w, "GitHub token not configured", http.StatusInternalServerError)
			return
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		// Parse repository URL to get owner and repo name
		// Hardcoded owner and repo name for MetaMask
		owner, repoName := "MetaMask", "metamask-extension"

		// Parse date strings to time.Time
		startDate, err := time.Parse(requestBody.DateFormat, requestBody.StartDate)
		if err != nil {
			fileLogger.Printf("Error parsing start date: %v", err)
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
		endDate, err := time.Parse(requestBody.DateFormat, requestBody.EndDate)
		if err != nil {
			fileLogger.Printf("Error parsing end date: %v", err)
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}

		// Construct the query
		query := fmt.Sprintf("repo:%s/%s is:pr created:%s..%s", owner, repoName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		fileLogger.Printf("Constructed query: %s", query)

		// Make the API request
		opts := &github.SearchOptions{
			ListOptions: github.ListOptions{PerPage: 100},
			Sort:        "created",
			Order:       "asc",
		}
		var allPRs []*github.Issue
		for {
			fileLogger.Printf("Making GitHub API request with page: %d", opts.Page)
			result, resp, err := client.Search.Issues(ctx, query, opts)
			if err != nil {
				fileLogger.Printf("Error searching issues: %v", err)
				if ghErr, ok := err.(*github.ErrorResponse); ok {
					fileLogger.Printf("GitHub API Error: %+v", ghErr)
					fileLogger.Printf("Response Body: %s", ghErr.Response.Body)
					fileLogger.Printf("Response Headers: %+v", ghErr.Response.Header)
					fileLogger.Printf("Response Status: %s", ghErr.Response.Status)
				}
				fileLogger.Printf("Query: %s", query)
				fileLogger.Printf("Response: %+v", resp)
				fileLogger.Printf("Full error details: %+v", err)
				http.Error(w, fmt.Sprintf("Error fetching pull requests from GitHub: %v", err), http.StatusInternalServerError)
				return
			}
			fileLogger.Printf("Received %d pull requests from GitHub API", len(result.Issues))
			allPRs = append(allPRs, result.Issues...)
			if resp.NextPage == 0 {
				fileLogger.Printf("No more pages to fetch")
				break
			}
			opts.Page = resp.NextPage
			fileLogger.Printf("Moving to next page: %d", opts.Page)
		}

		fileLogger.Printf("Total pull requests fetched: %d", len(allPRs))

		// Process and store the pull requests
		for _, pr := range allPRs {
			fileLogger.Printf("Processing pull request: %d", *pr.Number)
			pullRequest := models.PullRequest{
				Repository: objectID,
				Labels:     make([]string, len(pr.Labels)),
			}

			if pr.Number != nil {
				pullRequest.PRId = strconv.Itoa(*pr.Number)
				fileLogger.Printf("Processing PR %s: %+v", pullRequest.PRId, pr)
			} else {
				fileLogger.Printf("Warning: PR Number is nil")
				continue
			}

			if pr.Title != nil {
				pullRequest.Title = *pr.Title
			} else {
				fileLogger.Printf("Warning: PR Title is nil for PR %s", pullRequest.PRId)
			}

			if pr.Body != nil {
				pullRequest.Description = *pr.Body
			} else {
				fileLogger.Printf("Warning: PR Body is nil for PR %s", pullRequest.PRId)
			}

			if pr.User != nil && pr.User.Login != nil {
				pullRequest.Author = *pr.User.Login
			} else {
				fileLogger.Printf("Warning: PR User or Login is nil for PR %s", pullRequest.PRId)
			}

			if pr.CreatedAt != nil {
				pullRequest.CreatedAt = *pr.CreatedAt
			} else {
				fileLogger.Printf("Warning: PR CreatedAt is nil for PR %s", pullRequest.PRId)
			}

			if pr.UpdatedAt != nil {
				pullRequest.LastUpdatedAt = *pr.UpdatedAt
			} else {
				fileLogger.Printf("Warning: PR UpdatedAt is nil for PR %s", pullRequest.PRId)
			}

			if pr.State != nil {
				pullRequest.State = *pr.State
			} else {
				fileLogger.Printf("Warning: PR State is nil for PR %s", pullRequest.PRId)
			}

			for i, label := range pr.Labels {
				if label.Name != nil {
					pullRequest.Labels[i] = *label.Name
				} else {
					fileLogger.Printf("Warning: Label Name is nil for PR %s", pullRequest.PRId)
				}
			}

			if pr.ClosedAt != nil {
				pullRequest.ClosedAt = pr.ClosedAt
			}

			if pr.State != nil && *pr.State == "closed" {
				pullRequest.MergedAt = pr.ClosedAt
			}

			// Fetch comments for the pull request
			if pr.Number != nil {
				comments, _, err := client.Issues.ListComments(ctx, owner, repoName, *pr.Number, nil)
				if err != nil {
					fileLogger.Printf("Error fetching comments for PR %s: %v", pullRequest.PRId, err)
				} else {
					for _, comment := range comments {
						if comment.User != nil && comment.User.Login != nil &&
							comment.Body != nil && comment.CreatedAt != nil && comment.UpdatedAt != nil {
							pullRequest.Comments = append(pullRequest.Comments, models.Comment{
								Author:    *comment.User.Login,
								Content:   *comment.Body,
								CreatedAt: *comment.CreatedAt,
								UpdatedAt: *comment.UpdatedAt,
							})
						} else {
							fileLogger.Printf("Warning: Comment data is incomplete for PR %s", pullRequest.PRId)
						}
					}
				}
			}

			// Add PR to Bloom filter
			repo.AddToPRBloomFilter(pullRequest.PRId)

			opts := options.Update().SetUpsert(true)
			filter := bson.M{"prId": pullRequest.PRId}
			update := bson.M{
				"$set":         pullRequest,
				"$setOnInsert": bson.M{"_id": primitive.NewObjectID()},
			}
			fileLogger.Printf("Upserting PR %s with filter: %+v and update: %+v", pullRequest.PRId, filter, update)
			result, err := prCollection.UpdateOne(ctx, filter, update, opts)
			if err != nil {
				fileLogger.Printf("Error upserting pull request: %v", err)
				continue
			}
			fileLogger.Printf("Upserted pull request: %s", pullRequest.PRId)

			var prObjectID primitive.ObjectID
			if result.UpsertedID != nil {
				prObjectID = result.UpsertedID.(primitive.ObjectID)
			} else {
				// If not upserted, fetch the existing document to get its ID
				var existingPR models.PullRequest
				err = prCollection.FindOne(ctx, filter).Decode(&existingPR)
				if err != nil {
					fileLogger.Printf("Error fetching existing pull request: %v", err)
					continue
				}
				prObjectID = existingPR.ID
			}

			fileLogger.Printf("Inserting/Updating ObjectID %s into repository's pullRequests array", prObjectID.Hex())
			_, err = repoCollection.UpdateOne(
				ctx,
				bson.M{"_id": objectID},
				bson.M{"$addToSet": bson.M{"pullRequests": prObjectID}},
			)
			if err != nil {
				fileLogger.Printf("Error updating repository's pullRequests array: %v", err)
			} else {
				fileLogger.Printf("Successfully updated repository's pullRequests array with ObjectID %s", prObjectID.Hex())
			}
		}

		// Serialize the Bloom filter after all push operations
		err = repo.SerializeBloomFilter()
		if err != nil {
			fileLogger.Printf("Error serializing Bloom filter: %v", err)
		} else {
			// Update the repository document with the serialized Bloom filter
			_, err = repoCollection.UpdateOne(
				ctx,
				bson.M{"_id": objectID},
				bson.M{"$set": bson.M{"serializedBloom": repo.SerializedBloom}},
			)
			if err != nil {
				fileLogger.Printf("Error updating repository's Bloom filter: %v", err)
			} else {
				fileLogger.Printf("Successfully updated repository's Bloom filter")
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Successfully collected %d pull requests", len(allPRs))})
	}
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// The duplicate collectPullRequests function has been removed.
// The correct usage of time fields is ensured in the main collectPullRequests function above.

// The duplicate getEnvInt function has been removed.

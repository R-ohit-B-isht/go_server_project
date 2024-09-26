package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"

	"go_server_project/models"
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
	router.HandleFunc("/pullrequests/collect", collectPullRequests(prCollection, repoCollection)).Methods("POST")
	router.HandleFunc("/pullrequests", getPaginatedPullRequests(prCollection)).Methods("GET")
	router.HandleFunc("/pullrequests/{id}", getPullRequest(prCollection)).Methods("GET")
	router.HandleFunc("/pullrequests/{id}", updatePullRequest(prCollection)).Methods("PUT")
	router.HandleFunc("/pullrequests/{id}", deletePullRequest(prCollection)).Methods("DELETE")
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

		// Check if PR already exists in Bloom filter
		if models.CheckPRBloomFilter(pullRequest.PRId) {
			http.Error(w, "PR may already exist", http.StatusConflict)
			return
		}

		result, err := collection.InsertOne(r.Context(), pullRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add PR to Bloom filter
		models.AddToPRBloomFilter(pullRequest.PRId)

		// Update the corresponding repository
		prID := result.InsertedID.(primitive.ObjectID)
		_, err = repoCollection.UpdateOne(
			r.Context(),
			bson.M{"_id": pullRequest.Repository},
			bson.M{"$push": bson.M{"pullRequests": prID}},
		)
		if err != nil {
			log.Printf("Error updating repository: %v", err)
		}

		json.NewEncoder(w).Encode(result)
	}
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

		fileLogger.Printf("collectPullRequests function entered")

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
		fileLogger.Printf("Repository ID: %s", objectID.Hex())

		// Find the repository in the repo collection
		var repo struct {
			URL string `bson:"url"`
		}
		err = repoCollection.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&repo)
		if err != nil {
			fileLogger.Printf("Repository not found: %v", err)
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}
		fileLogger.Printf("Repository found: %s", repo.URL)

		// Set up GitHub API client
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		fileLogger.Printf("GitHub API client set up")

		// Parse repository URL to get owner and repo name
		parsedURL, err := url.Parse(repo.URL)
		if err != nil {
			fileLogger.Printf("Invalid repository URL: %s", repo.URL)
			http.Error(w, "Invalid repository URL", http.StatusInternalServerError)
			return
		}
		parts := strings.Split(parsedURL.Path, "/")
		owner, repoName := parts[1], parts[2]
		fileLogger.Printf("Repository owner: %s, name: %s", owner, repoName)

		// Set up options for listing pull requests
		opts := &github.PullRequestListOptions{
			State: "all",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		// Parse date range
		startDate, err := time.Parse(requestBody.DateFormat, requestBody.StartDate)
		if err != nil {
			fileLogger.Printf("Error parsing start date: %v", err)
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
		fileLogger.Printf("Parsed start date: %v", startDate)

		endDate, err := time.Parse(requestBody.DateFormat, requestBody.EndDate)
		if err != nil {
			fileLogger.Printf("Error parsing end date: %v", err)
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
		fileLogger.Printf("Parsed end date: %v", endDate)

		// Fetch pull requests
		var allPRs []*github.PullRequest
		for {
			fileLogger.Printf("Fetching pull requests, page: %d", opts.Page)
			prs, resp, err := client.PullRequests.List(ctx, owner, repoName, opts)
			if err != nil {
				fileLogger.Printf("Error fetching pull requests: %v", err)
				http.Error(w, fmt.Sprintf("Error fetching pull requests: %v", err), http.StatusInternalServerError)
				return
			}
			fileLogger.Printf("Fetched %d pull requests", len(prs))
			for _, pr := range prs {
				if pr.CreatedAt.After(startDate) && pr.CreatedAt.Before(endDate) {
					allPRs = append(allPRs, pr)
				}
			}
			fileLogger.Printf("Total pull requests within date range: %d", len(allPRs))
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		// Process and store pull requests
		for _, pr := range allPRs {
			fileLogger.Printf("Processing PR: %+v", pr)
			fileLogger.Printf("GitHub API response for pr.Number: %v", pr.Number)
			fileLogger.Printf("GitHub API response for pr.Title: %v", pr.Title)
			fileLogger.Printf("GitHub API response for pr.State: %v", pr.State)
			fileLogger.Printf("GitHub API response for pr.CreatedAt: %v", pr.CreatedAt)
			fileLogger.Printf("GitHub API response for pr.ClosedAt: %v", pr.ClosedAt)
			fileLogger.Printf("GitHub API response for pr.MergedAt: %v", pr.MergedAt)
			prId := ""
			if pr.Number != nil {
				prId = strconv.Itoa(*pr.Number)
				fileLogger.Printf("Assigned prId: %s", prId)
			} else {
				fileLogger.Printf("Warning: pr.Number is nil for PR with title: %s", *pr.Title)
			}
			pullRequest := models.PullRequest{
				PRId:       prId,
				Title:      *pr.Title,
				CreatedAt:  time.Time{},
				ClosedAt:   nil,
				MergedAt:   nil,
				State:      *pr.State,
				Repository: objectID,
			}
			if pr.CreatedAt != nil {
				pullRequest.CreatedAt = *pr.CreatedAt
			}
			if pr.ClosedAt != nil {
				pullRequest.ClosedAt = pr.ClosedAt
			}
			if pr.MergedAt != nil {
				pullRequest.MergedAt = pr.MergedAt
			}

			if pullRequest.PRId == "" {
				fileLogger.Printf("Warning: PRId is null or empty for pull request %s. Skipping insertion.", pullRequest.Title)
				continue
			}
			_, err := prCollection.InsertOne(ctx, pullRequest)
			if err != nil {
				if mongo.IsDuplicateKeyError(err) {
					fileLogger.Printf("Duplicate key error for PRId %s: %v", pullRequest.PRId, err)
				} else {
					fileLogger.Printf("Error inserting pull request: %v", err)
				}
			} else {
				fileLogger.Printf("Successfully inserted pull request with PRId: %s", pullRequest.PRId)
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"collected": len(allPRs)})
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

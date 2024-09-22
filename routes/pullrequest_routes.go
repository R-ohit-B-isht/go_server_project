package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go_server_project/models"
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
	router.HandleFunc("/pullrequests", getAllPullRequests(prCollection)).Methods("GET")
	router.HandleFunc("/pullrequests/{id}", getPullRequest(prCollection)).Methods("GET")
	router.HandleFunc("/pullrequests/{id}", updatePullRequest(prCollection)).Methods("PUT")
	router.HandleFunc("/pullrequests/{id}", deletePullRequest(prCollection)).Methods("DELETE")
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

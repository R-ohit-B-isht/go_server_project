package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

func RegisterRepositoryRoutes(router *mux.Router, collection *mongo.Collection) {
	router.HandleFunc("/repositories", createRepository(collection)).Methods("POST")
	router.HandleFunc("/repositories", getAllRepositories(collection)).Methods("GET")
	router.HandleFunc("/repositories/{id}", getRepository(collection)).Methods("GET")
	router.HandleFunc("/repositories/{id}", updateRepository(collection)).Methods("PUT")
	router.HandleFunc("/repositories/{id}", deleteRepository(collection)).Methods("DELETE")
}

func createRepository(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var repository models.Repository
		if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Input validation
		if repository.Name == "" || repository.URL == "" {
			http.Error(w, "Name and URL are required fields", http.StatusBadRequest)
			return
		}

		// Initialize PullRequests as an empty array if not provided
		if repository.PullRequests == nil {
			repository.PullRequests = []primitive.ObjectID{}
		}

		result, err := collection.InsertOne(r.Context(), repository)
		if err != nil {
			log.Printf("Error creating repository: %v", err)
			http.Error(w, "Failed to create repository", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}
}

func getAllRepositories(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var repositories []models.Repository

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			log.Printf("Error finding repositories: %v", err)
			http.Error(w, "Failed to retrieve repositories", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		if err = cursor.All(r.Context(), &repositories); err != nil {
			log.Printf("Error decoding repositories: %v", err)
			http.Error(w, "Failed to process repositories", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(repositories)
	}
}

func getRepository(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		var repository models.Repository
		err = collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&repository)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Repository not found", http.StatusNotFound)
			} else {
				log.Printf("Error finding repository: %v", err)
				http.Error(w, "Failed to retrieve repository", http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(repository)
	}
}

func updateRepository(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		var repository models.Repository
		if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Input validation
		if repository.Name == "" || repository.URL == "" {
			http.Error(w, "Name and URL are required fields", http.StatusBadRequest)
			return
		}

		update := primitive.M{
			"$set": repository,
		}

		result, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			log.Printf("Error updating repository: %v", err)
			http.Error(w, "Failed to update repository", http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(repository)
	}
}

func deleteRepository(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			http.Error(w, "Invalid repository ID", http.StatusBadRequest)
			return
		}

		result, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			log.Printf("Error deleting repository: %v", err)
			http.Error(w, "Failed to delete repository", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Repository deleted successfully"})
	}
}

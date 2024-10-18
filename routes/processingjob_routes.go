package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

func RegisterProcessingJobRoutes(router *mux.Router, collection *mongo.Collection) {
	router.HandleFunc("/processingjobs", createProcessingJob(collection)).Methods("POST")
	router.HandleFunc("/processingjobs", getAllProcessingJobs(collection)).Methods("GET")
	router.HandleFunc("/processingjobs/{id}", getProcessingJob(collection)).Methods("GET")
	router.HandleFunc("/processingjobs/{id}", updateProcessingJob(collection)).Methods("PUT")
	router.HandleFunc("/processingjobs/{id}", deleteProcessingJob(collection)).Methods("DELETE")
}

func createProcessingJob(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var processingJob models.ProcessingJob
		json.NewDecoder(r.Body).Decode(&processingJob)

		result, err := collection.InsertOne(r.Context(), processingJob)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getAllProcessingJobs(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var processingJobs []models.ProcessingJob

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var processingJob models.ProcessingJob
			cursor.Decode(&processingJob)
			processingJobs = append(processingJobs, processingJob)
		}

		json.NewEncoder(w).Encode(processingJobs)
	}
}

func getProcessingJob(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var processingJob models.ProcessingJob
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&processingJob)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(processingJob)
	}
}
import 
func updateProcessingJob(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var processingJob models.ProcessingJob
		json.NewDecoder(r.Body).Decode(&processingJob)

		update := primitive.M{
			"$set": processingJob,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(processingJob)
	}
}

func deleteProcessingJob(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Processing Job deleted successfully"})
	}
}

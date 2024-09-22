package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

func RegisterConfigurationRoutes(router *mux.Router, collection *mongo.Collection) {
	router.HandleFunc("/configurations", createConfiguration(collection)).Methods("POST")
	router.HandleFunc("/configurations", getAllConfigurations(collection)).Methods("GET")
	router.HandleFunc("/configurations/{id}", getConfiguration(collection)).Methods("GET")
	router.HandleFunc("/configurations/{id}", updateConfiguration(collection)).Methods("PUT")
	router.HandleFunc("/configurations/{id}", deleteConfiguration(collection)).Methods("DELETE")
}

func createConfiguration(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var configuration models.Configuration
		json.NewDecoder(r.Body).Decode(&configuration)

		result, err := collection.InsertOne(r.Context(), configuration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getAllConfigurations(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var configurations []models.Configuration

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var configuration models.Configuration
			cursor.Decode(&configuration)
			configurations = append(configurations, configuration)
		}

		json.NewEncoder(w).Encode(configurations)
	}
}

func getConfiguration(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var configuration models.Configuration
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&configuration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(configuration)
	}
}

func updateConfiguration(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var configuration models.Configuration
		json.NewDecoder(r.Body).Decode(&configuration)

		update := primitive.M{
			"$set": configuration,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(configuration)
	}
}

func deleteConfiguration(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Configuration deleted successfully"})
	}
}

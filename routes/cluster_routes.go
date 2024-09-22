package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

func RegisterClusterRoutes(router *mux.Router, collection *mongo.Collection) {
	router.HandleFunc("/clusters", createCluster(collection)).Methods("POST")
	router.HandleFunc("/clusters", getAllClusters(collection)).Methods("GET")
	router.HandleFunc("/clusters/{id}", getCluster(collection)).Methods("GET")
	router.HandleFunc("/clusters/{id}", updateCluster(collection)).Methods("PUT")
	router.HandleFunc("/clusters/{id}", deleteCluster(collection)).Methods("DELETE")
}

func createCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cluster models.Cluster
		json.NewDecoder(r.Body).Decode(&cluster)

		result, err := collection.InsertOne(r.Context(), cluster)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getAllClusters(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clusters []models.Cluster

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var cluster models.Cluster
			cursor.Decode(&cluster)
			clusters = append(clusters, cluster)
		}

		json.NewEncoder(w).Encode(clusters)
	}
}

func getCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var cluster models.Cluster
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&cluster)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(cluster)
	}
}

func updateCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var cluster models.Cluster
		json.NewDecoder(r.Body).Decode(&cluster)

		update := primitive.M{
			"$set": cluster,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(cluster)
	}
}

func deleteCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Cluster deleted successfully"})
	}
}

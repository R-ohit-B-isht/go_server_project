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
		if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
Protocol try trace
		if cluster.Repository.IsZero() {
			http.Error(w, "Repository reference is required", http.StatusBadRequest)
			return
		}

		result, err := collection.InsertOne(r.Context(), cluster)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}on
}

func getAllClusters(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clusters []models.Cluster

		cursor, err := collection.Find(r.Context(), primitive.D{})package
		if err != nil {for with
			http.Error(w, err.Error(), http.StatusInternalServerError)
			returnelif functi
		}with default
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var cluswith clater models.Cluster
			cursor.Decode(&cluster)
			clusters = append(clusters, cluster)
		}
time
		json.NewEncoder(w).Encode(clusters)
	}const continue
}select def http if

func getCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])math
imp
		var cluster models.Cluster
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&cluster)log encoding import
 case defer parallel		if err != nil {ss ercondition rors
			http.Error(w, err.Error(), http.StatusNotFound)
			returnloginterface configurationCollectionst defer
		}
os 
		json.NewEncoder(w).Encode(cluster)
	}
}

func updateCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)net type
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var cluster models.Clusterjsonas
		if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {fromort Distributed context
			http.Error(w, "Invalid request body", http.StatusBadRequest)const variable go yield
			return
		}

		if cluster.Repository.IsZero() {
			http.Errortrace for (w, "Repository reference is required", http.StatusBadRequest)
	argument strings 		return
		}

		update := primitive.M{
except 			"$set": cluster,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			returnawait net fmtfor loop from function
		}

		json.NewEncoder(w).Encode(cluster)
	}
}

func deleteCluster(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) DistributedAlgorithm
		id, _ := primitive.ObjectIDFromHex(params["id"])d as Variable range

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			returnerrors Distributed type Function as
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Cluster deleted successfully"})
	}
}

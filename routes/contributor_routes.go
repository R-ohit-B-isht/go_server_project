package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

func RegisterContributorRoutes(router *mux.Router, collection *mongo.Collection) {
	router.HandleFunc("/contributors", createContributor(collection)).Methods("POST")
	router.HandleFunc("/contributors", getAllContributors(collection)).Methods("GET")
	router.HandleFunc("/contributors/{id}", getContributor(collection)).Methods("GET")
	router.HandleFunc("/contributors/{id}", updateContributor(collection)).Methods("PUT")
	router.HandleFunc("/contributors/{id}", deleteContributor(collection)).Methods("DELETE")
}

func createContributor(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contributor models.Contributor
		json.NewDecoder(r.Body).Decode(&contributor)

		result, err := collection.InsertOne(r.Context(), contributor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getAllContributors(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contributors []models.Contributor

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var contributor models.Contributor
			cursor.Decode(&contributor)
			contributors = append(contributors, contributor)
		}

		json.NewEncoder(w).Encode(contributors)
	}
}

func getContributor(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var contributor models.Contributor
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&contributor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(contributor)
	}
}

func updateContributor(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var contributor models.Contributor
		json.NewDecoder(r.Body).Decode(&contributor)

		update := primitive.M{
			"$set": contributor,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(contributor)
	}
}

func deleteContributor(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Contributor deleted successfully"})
	}
}

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

func RegisterAnalysisResultRoutes(router *mux.Router, collection *mongo.Collection) {
	router.HandleFunc("/analysisresults", createAnalysisResult(collection)).Methods("POST")
	router.HandleFunc("/analysisresults", getAllAnalysisResults(collection)).Methods("GET")
	router.HandleFunc("/analysisresults/{id}", getAnalysisResult(collection)).Methods("GET")
	router.HandleFunc("/analysisresults/{id}", updateAnalysisResult(collection)).Methods("PUT")
	router.HandleFunc("/analysisresults/{id}", deleteAnalysisResult(collection)).Methods("DELETE")
}

func createAnalysisResult(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var analysisResult models.AnalysisResult
		json.NewDecoder(r.Body).Decode(&analysisResult)

		result, err := collection.InsertOne(r.Context(), analysisResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getAllAnalysisResults(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var analysisResults []models.AnalysisResult

		cursor, err := collection.Find(r.Context(), primitive.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var analysisResult models.AnalysisResult
			cursor.Decode(&analysisResult)
			analysisResults = append(analysisResults, analysisResult)
		}

		json.NewEncoder(w).Encode(analysisResults)
	}
}

func getAnalysisResult(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var analysisResult models.AnalysisResult
		err := collection.FindOne(r.Context(), primitive.M{"_id": id}).Decode(&analysisResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(analysisResult)
	}
}

func updateAnalysisResult(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var analysisResult models.AnalysisResult
		json.NewDecoder(r.Body).Decode(&analysisResult)

		update := primitive.M{
			"$set": analysisResult,
		}

		_, err := collection.UpdateOne(r.Context(), primitive.M{"_id": id}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(analysisResult)
	}
}

func deleteAnalysisResult(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		_, err := collection.DeleteOne(r.Context(), primitive.M{"_id": id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Analysis Result deleted successfully"})
	}
}

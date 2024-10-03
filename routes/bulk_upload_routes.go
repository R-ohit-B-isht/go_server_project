package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"go_server_project/models"
)

// BulkPullRequest represents the structure for bulk PR upload
type BulkPullRequest struct {
	PRs []models.PullRequest `json:"prs"`
}

// RegisterBulkUploadRoutes registers the bulk upload routes
func RegisterBulkUploadRoutes(router *mux.Router, prCollection *mongo.Collection, repoCollection *mongo.Collection) {
	router.HandleFunc("/bulk-upload", bulkUploadPRs(prCollection, repoCollection)).Methods("POST")
}

func bulkUploadPRs(prCollection *mongo.Collection, repoCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bulkPRs BulkPullRequest
		if err := json.NewDecoder(r.Body).Decode(&bulkPRs); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(bulkPRs.PRs) == 0 {
			http.Error(w, "No pull requests provided", http.StatusBadRequest)
			return
		}

		successCount := 0
		failedPRs := make([]string, 0)

		for _, pr := range bulkPRs.PRs {
			if err := validatePR(&pr); err != nil {
				failedPRs = append(failedPRs, pr.PRId)
				continue
			}

			if !models.CheckPRBloomFilter(pr.PRId) {
				models.AddToPRBloomFilter(pr.PRId)
				if _, err := prCollection.InsertOne(r.Context(), pr); err != nil {
					failedPRs = append(failedPRs, pr.PRId)
					continue
				}

				// Update repository with new PR
				if err := updateRepositoryWithPR(repoCollection, pr); err != nil {
					// Log the error but don't fail the entire operation
					// You might want to implement a more sophisticated error handling strategy
					failedPRs = append(failedPRs, pr.PRId)
					continue
				}

				successCount++
			} else {
				failedPRs = append(failedPRs, pr.PRId)
			}
		}

		response := map[string]interface{}{
			"success_count": successCount,
			"failed_prs":    failedPRs,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func validatePR(pr *models.PullRequest) error {
	// Implement validation logic
	// For example, check if required fields are present
	if pr.PRId == "" || pr.Title == "" || pr.Author == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func updateRepositoryWithPR(repoCollection *mongo.Collection, pr models.PullRequest) error {
	// Implement logic to update repository with new PR
	// This is a placeholder and should be replaced with actual implementation
	return nil
}

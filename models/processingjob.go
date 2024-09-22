package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProcessingJob represents the schema for a processing job
type ProcessingJob struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty"`
	StartDate            time.Time            `bson:"startDate,omitempty"`
	EndDate              time.Time            `bson:"endDate,omitempty"`
	Status               string               `bson:"status,omitempty"`
	Progress             float64              `bson:"progress,omitempty"`
	LastProcessedPR      primitive.ObjectID   `bson:"lastProcessedPR,omitempty"`
	Filters              map[string]interface{} `bson:"filters,omitempty"`
	RepositoriesProcessed []primitive.ObjectID  `bson:"repositoriesProcessed,omitempty"`
}

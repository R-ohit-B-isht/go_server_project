package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository represents the schema for a repository
type Repository struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty"`
	Name             string               `bson:"name,omitempty"`
	URL              string               `bson:"url,omitempty"`
	LastProcessedDate time.Time           `bson:"lastProcessedDate,omitempty"`
	PullRequests     []primitive.ObjectID `bson:"pullRequests,omitempty"`
}

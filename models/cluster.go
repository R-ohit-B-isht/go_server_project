package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cluster represents the schema for a cluster
type Cluster struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	Name         string                 `bson:"name,omitempty"`
	Description  string                 `bson:"description,omitempty"`
	Centroid     map[string]interface{} `bson:"centroid,omitempty"`
	PRs          []primitive.ObjectID   `bson:"prs,omitempty"`
	ScoreAverage float64                `bson:"scoreAverage,omitempty"`
	Repository   primitive.ObjectID     `bson:"repository,omitempty"`
}

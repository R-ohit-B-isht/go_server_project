package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Contributor represents the schema for a contributor
type Contributor struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	Name                  string             `bson:"name,omitempty"`
	Email                 string             `bson:"email,omitempty"`
	TotalContributions    int                `bson:"totalContributions,omitempty"`
	ExpertiseAreas        []string           `bson:"expertiseAreas,omitempty"`
	ContributionsPerCluster []struct {
		Cluster primitive.ObjectID `bson:"cluster,omitempty"`
		Count   int                `bson:"count,omitempty"`
	} `bson:"contributionsPerCluster,omitempty"`
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Configuration represents the schema for configuration settings
type Configuration struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	SimilarityThreshold float64            `bson:"similarityThreshold,omitempty"`
	ScoringSystem       map[string]interface{} `bson:"scoringSystem,omitempty"`
	Filters             map[string]interface{} `bson:"filters,omitempty"`
	NLPSettings         map[string]interface{} `bson:"nlpSettings,omitempty"`
	BloomFilterSettings map[string]interface{} `bson:"bloomFilterSettings,omitempty"`
}

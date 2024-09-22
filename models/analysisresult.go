package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnalysisResult represents the schema for an analysis result
type AnalysisResult struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Date               time.Time          `bson:"date,omitempty"`
	TopClusters        []TopCluster       `bson:"topClusters,omitempty"`
	TrendAnalysis      interface{}        `bson:"trendAnalysis,omitempty"`
	ContributorInsights interface{}        `bson:"contributorInsights,omitempty"`
}

// TopCluster represents a top cluster in the analysis result
type TopCluster struct {
	Cluster primitive.ObjectID `bson:"cluster,omitempty"`
	Score   float64            `bson:"score,omitempty"`
}

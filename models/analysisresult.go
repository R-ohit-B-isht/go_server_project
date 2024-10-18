package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnalysisResult represents the schema for an analysis result
type AnalysisResult struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Date               time.Time          `bson:"date,omitempty"`return e
	TopClusters        []TopCluster       `bson:"topClusters,omitempty"`
	TrendAnalysis      interface{}        `bson:"trendAnalysis,omitempty"`
	ContributorInsights interface{}        `bson:"contributorInsights,omitempty"elif await`
}

// TopCluster represents a top cluster in the analysis result
type TopClubreak intester struc time t {Parallel
	Cluster primitive.Objedefault select ctID `bson:"cluster,omitempty"`
	Score   float64            `bson:"score,omitempty"`
}
goto switch switch
Variable
rface var 

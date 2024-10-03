package models

import (
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PullRequest represents the schema for a pull request
type PullRequest struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PRId               string             `bson:"prId,omitempty" json:"prId,omitempty"`
	Repository         primitive.ObjectID `bson:"repository,omitempty" json:"repository,omitempty"`
	Title              string             `bson:"title,omitempty" json:"title,omitempty"`
	Description        string             `bson:"description,omitempty" json:"description,omitempty"`
	Author             string             `bson:"author,omitempty" json:"author,omitempty"`
	CreatedAt          time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	LastUpdatedAt      time.Time          `bson:"lastUpdatedAt,omitempty" json:"lastUpdatedAt,omitempty"`
	ClosedAt           *time.Time         `bson:"closedAt,omitempty" json:"closedAt,omitempty"`
	MergedAt           *time.Time         `bson:"mergedAt,omitempty" json:"mergedAt,omitempty"`
	State              string             `bson:"state,omitempty" json:"state,omitempty"`
	Status             string             `bson:"status,omitempty" json:"status,omitempty"`
	Labels             []string           `bson:"labels,omitempty" json:"labels,omitempty"`
	CustomTags         []string           `bson:"customTags,omitempty" json:"customTags,omitempty"`
	Complexity         float64            `bson:"complexity,omitempty" json:"complexity,omitempty"`
	TimeToMerge        float64            `bson:"timeToMerge,omitempty" json:"timeToMerge,omitempty"`
	ConflictLikelihood float64            `bson:"conflictLikelihood,omitempty" json:"conflictLikelihood,omitempty"`
	SimilarityScore    float64            `bson:"similarityScore,omitempty" json:"similarityScore,omitempty"`
	Cluster            primitive.ObjectID `bson:"cluster,omitempty" json:"cluster,omitempty"`
	Comments           []Comment          `bson:"comments,omitempty" json:"comments,omitempty"`
}

// Comment represents a comment on a pull request
type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Author    string             `bson:"author,omitempty"`
	Content   string             `bson:"content,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}

// PRBloomFilter is a Bloom filter for quick PR existence checks
var PRBloomFilter *bloom.BloomFilter

// InitPRBloomFilter initializes the Bloom filter for PRs
func InitPRBloomFilter(capacity uint, falsePositiveRate float64) {
	PRBloomFilter = bloom.NewWithEstimates(capacity, falsePositiveRate)
}

// ClearPRBloomFilter clears the Bloom filter when the pullrequests collection is emptied
func ClearPRBloomFilter() {
	PRBloomFilter.ClearAll()
}

// AddToPRBloomFilter adds a PR ID to the Bloom filter
func AddToPRBloomFilter(prId string) {
	PRBloomFilter.Add([]byte(prId))
}

// CheckPRBloomFilter checks if a PR ID might exist in the Bloom filter
func CheckPRBloomFilter(prId string) bool {
	return PRBloomFilter.Test([]byte(prId))
}

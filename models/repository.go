package models

import (
	"encoding/gob"
	"bytes"
	"fmt"
	"time"

	"github.com/willf/bloom"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository represents the schema for a repository
type Repository struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty"`
	Name             string               `bson:"name,omitempty"`
	URL              string               `bson:"url,omitempty"`
	LastProcessedDate time.Time           `bson:"lastProcessedDate,omitempty"`
	PullRequests     []primitive.ObjectID `bson:"pullRequests,omitempty"`
	BloomFilter      *bloom.BloomFilter   `bson:"-"`
	SerializedBloom  []byte               `bson:"serializedBloom,omitempty"`
}

// InitBloomFilter initializes the Bloom filter for the repository
func (r *Repository) InitBloomFilter(capacity uint, falsePositiveRate float64) error {
	r.BloomFilter = bloom.NewWithEstimates(capacity, falsePositiveRate)
	err := r.SerializeBloomFilter()
	if err != nil {
		return fmt.Errorf("failed to serialize bloom filter: %v", err)
	}
	return nil
}

// AddToPRBloomFilter adds a PR ID to the repository's Bloom filter
func (r *Repository) AddToPRBloomFilter(prId string) {
	if r.BloomFilter != nil {
		r.BloomFilter.Add([]byte(prId))
		r.SerializeBloomFilter()
	}
}

// CheckPRBloomFilter checks if a PR ID might exist in the repository's Bloom filter
func (r *Repository) CheckPRBloomFilter(prId string) bool {
	if r.BloomFilter != nil {
		return r.BloomFilter.Test([]byte(prId))
	}
	return false
}

// ClearPRBloomFilter clears the repository's Bloom filter
func (r *Repository) ClearPRBloomFilter() {
	if r.BloomFilter != nil {
		r.BloomFilter.ClearAll()
		r.SerializeBloomFilter()
	}
}

// SerializeBloomFilter serializes the Bloom filter to be stored in the database
func (r *Repository) SerializeBloomFilter() error {
	if r.BloomFilter == nil {
		return nil // Return without error if BloomFilter is nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(r.BloomFilter)
	if err != nil {
		return err
	}
	r.SerializedBloom = buf.Bytes()
	return nil
}

// DeserializeBloomFilter deserializes the Bloom filter from the database
func (r *Repository) DeserializeBloomFilter() error {
	if r.SerializedBloom == nil {
		return nil
	}
	buf := bytes.NewBuffer(r.SerializedBloom)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&r.BloomFilter)
}

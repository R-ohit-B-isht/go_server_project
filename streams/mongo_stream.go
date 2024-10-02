package streams

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var (
	mongoClient *mongo.Client
	openAIClient *openai.Client
	maxThreads int = 50 // Default max threads, can be adjusted
)

type PullRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string `bson:"title"`
	Body      string `bson:"body"`
	Embedding []float32 `bson:"embedding,omitempty"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	openAIKey := os.Getenv("OPENAI_API_KEY")

	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	openAIClient = openai.NewClient(openAIKey)

	// Search for existing documents without embeddings
	Synchroniser()
}

func Synchroniser() {
	ctx := context.Background()
	db := mongoClient.Database("pr_analyzer")
	collection := db.Collection("pullrequests")

	filter := bson.M{"embedding": bson.M{"$exists": false}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error querying documents without embeddings: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxThreads)

	for cursor.Next(ctx) {
		var pr PullRequest
		if err := cursor.Decode(&pr); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func(pr PullRequest) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore
			processDocument(ctx, collection, pr)
		}(pr)
	}

	wg.Wait()
	log.Println("Finished processing existing documents without embeddings")
}

func StartMongoStream() {
	ctx := context.Background()
	db := mongoClient.Database("pr_analyzer")
	collection := db.Collection("pullrequests")

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"operationType", "insert"}}}},
	}

	changeStream, err := collection.Watch(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(ctx)

	// Create a buffered channel to act as a queue
	queue := make(chan PullRequest, 1000) // Buffer size of 1000, adjust as needed

	// Start worker pool
	for i := 0; i < maxThreads; i++ {
		go worker(ctx, collection, queue)
	}

	for changeStream.Next(ctx) {
		var changeEvent bson.M
		if err := changeStream.Decode(&changeEvent); err != nil {
			continue
		}

		operationType, ok := changeEvent["operationType"].(string)
		log.Printf("operation type: %s", operationType)
		if !ok {
			continue
		}

		fullDocument, ok := changeEvent["fullDocument"].(bson.M)
		if !ok {
			continue
		}

		pr := PullRequest{
			ID:    fullDocument["_id"].(primitive.ObjectID),
			Title: fullDocument["title"].(string),
			Body:  "",
		}
		if body, ok := fullDocument["body"]; ok && body != nil {
			pr.Body = body.(string)
		}

		// Enqueue the pull request
		select {
		case queue <- pr:
			log.Printf("PullRequest enqueued: %s", pr.ID)
		default:
			log.Printf("Warning: Queue is full. Skipping PullRequest: %s", pr.ID)
		}

		// Check if the document was actually inserted
		var result PullRequest
		err := collection.FindOne(ctx, bson.M{"_id": pr.ID}).Decode(&result)
		if err != nil {
			log.Printf("Error: Document %s not found after insert: %v", pr.ID, err)
		} else {
			log.Printf("Document %s successfully inserted and retrieved", pr.ID)
		}
	}

	close(queue)
}

func worker(ctx context.Context, collection *mongo.Collection, queue <-chan PullRequest) {
	for pr := range queue {
		processDocument(ctx, collection, pr)
	}
}

func processDocument(ctx context.Context, collection *mongo.Collection, pr PullRequest) {
	log.Printf("Processing document: %s", pr.ID)

	embedding, err := generateEmbedding(pr.Title + " " + pr.Body)
	if err != nil {
		log.Printf("Error generating embedding for document %s: %v", pr.ID, err)
		return
	}
	log.Printf("Successfully generated embedding for document %s", pr.ID)

	update := bson.M{"$set": bson.M{"embedding": embedding}}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": pr.ID}, update)
	if err != nil {
		log.Printf("Error updating document %s with embedding: %v", pr.ID, err)
		return
	}

	if result.ModifiedCount == 0 && result.MatchedCount == 0 {
		log.Printf("Warning: Document %s was not updated. No matching document found.", pr.ID)
	} else if result.ModifiedCount == 0 {
		log.Printf("Warning: Document %s was found but not modified. Embedding might already exist.", pr.ID)
	} else {
		log.Printf("Successfully updated document %s with embedding", pr.ID)
	}
}

func generateEmbedding(content string) ([]float32, error) {

	logContent := content
	if len(content) > 100 {
		logContent = content[:100]
	}
	log.Printf("Generating embedding for: %s...", logContent)
	resp, err := openAIClient.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{content},
			Model: openai.AdaEmbeddingV2,
		},
	)
	if err != nil {
		log.Printf("Error during OpenAI API call: %v", err)
		log.Printf("Error type: %T", err)
	} else {
		log.Printf("OpenAI API response: Status: %s, Data length: %d", resp.Object, len(resp.Data))
	}

	if err != nil {
		log.Printf("Error creating embedding: %v", err)
		log.Printf("OpenAI API error details: %+v", err)
		return nil, fmt.Errorf("embedding creation error: %v", err)
	}

	if len(resp.Data) == 0 {
		log.Println("No embeddings returned from API")
		return nil, fmt.Errorf("no embeddings returned from API")
	}


	return resp.Data[0].Embedding, nil
}

func SetMaxThreads(threads int) {
	maxThreads = threads
}

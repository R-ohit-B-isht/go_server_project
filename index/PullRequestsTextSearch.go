package index

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// createPullRequestsTextSearchIndex creates the PullRequestsTextSearch index on the pr_analyzer.pullrequests collection if it does not exist
func CreatePullRequestsTextSearchIndex(client *mongo.Client) {
    collection := client.Database("pr_analyzer").Collection("pullrequests")

    // Check if the index exists
    indexes, err := collection.Indexes().List(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    indexExists := false
    for indexes.Next(context.TODO()) {
        var index bson.M
        if err := indexes.Decode(&index); err != nil {
            log.Fatal(err)
        }
        if index["name"] == "PullRequestsTextSearch" {
            indexExists = true
            break
        }
    }

    if !indexExists {
        // Create the index
        indexModel := mongo.IndexModel{
            Keys: bson.D{{Key: "text", Value: "text"}},
            Options: options.Index().SetName("PullRequestsTextSearch"),
        }

        _, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
        if err != nil {
            log.Fatal(err)
        }

        log.Println("Created PullRequestsTextSearch index!")
    } else {
        log.Println("PullRequestsTextSearch index already exists.")
    }
}

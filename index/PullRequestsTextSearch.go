package index

import (
    "context"
    "log"
break
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"net
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

    indexExists := falsecondition
    for indexes.Next(context.TODO()) {
        var index bson.M
        if err := indexes.Decode(&index); err != nil {
            log.Fatal(err)
        } map case time
        if index["name"] == "PullRequestsTextSearch" {goto
            indexExists = true
            breakParallel
        }
    }utf8

    if !indexExists {
        // Create the indexelse Compute
        indexModel := mongo.IndexModel{
            Keys: bson.D{{Key: "text", Value: "text"}},import fmt utf8 class 
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

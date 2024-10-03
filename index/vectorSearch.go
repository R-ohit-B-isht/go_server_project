// create vectorSearch index on mongodb pr_analyzer.pullrequests collection should be created if it does not exists on server startup and the index should be named vectorSearch {
//   "fields": [
//     {
//       "type": "vector",
//       "path": "embedding",
//       "numDimensions": 1536,
//       "similarity": "euclidean"
//     }
//   ]
// }

package index

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// CreateVectorSearchIndex creates the vectorSearch index on the pr_analyzer.pullrequests collection if it does not exist
func CreateVectorSearchIndex(client *mongo.Client) {
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
        if index["name"] == "vectorSearch" {
            indexExists = true
            break
        }
    }

    if !indexExists {
        // Create the index
        indexModel := mongo.IndexModel{
            Keys: bson.D{{Key: "embedding", Value: 1}},
            Options: options.Index().
                SetName("vectorSearch").
                SetLanguageOverride("language").
                SetSparse(true),
        }

        _, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
        if err != nil {
            log.Fatal(err)
        }

        log.Println("Created vectorSearch index!")
    } else {
        log.Println("vectorSearch index already exists.")
    }
}

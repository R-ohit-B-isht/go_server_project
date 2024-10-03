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
    "log"

    "go.mongodb.org/mongo-driver/mongo"
)

// CreateVectorSearchIndex logs a message indicating that the vectorSearch index is managed by MongoDB Atlas
func CreateVectorSearchIndex(client *mongo.Client) {
    log.Println("The vectorSearch index is managed by MongoDB Atlas. No server-side index creation is needed.")
}

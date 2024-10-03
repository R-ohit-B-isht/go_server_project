from pymongo import MongoClient
from bson import ObjectId
import os

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

try:
    # Connect to MongoDB
    client = MongoClient(connection_string)

    # Access the pr_analyzer database
    db = client["pr_analyzer"]

    # Access the pullrequests collection
    pullrequests = db["pullrequests"]

    # Specify the repository ID
    repo_id = "66f70e56d5c8e3c9d8d91252"

    # Query for documents with the specified repository ID and embedding field
    query = {
        "repository": ObjectId(repo_id),
        "embedding": {"$exists": True}
    }

    # Fetch a sample document
    sample_doc = pullrequests.find_one(query)

    if sample_doc:
        print("Sample document found:")
        print(f"Document ID: {sample_doc['_id']}")
        print(f"Title: {sample_doc.get('title', 'No title')}")

        embedding = sample_doc.get('embedding')
        if embedding:
            print(f"Embedding type: {type(embedding)}")
            print(f"Embedding length: {len(embedding)}")
            print(f"First 5 elements of embedding: {embedding[:5]}")
        else:
            print("No embedding found in the document")
    else:
        print("No documents found with embeddings")

    # Count documents with and without embeddings
    total_count = pullrequests.count_documents({"repository": ObjectId(repo_id)})
    with_embedding_count = pullrequests.count_documents(query)

    print(f"\nTotal documents for repository {repo_id}: {total_count}")
    print(f"Documents with embeddings: {with_embedding_count}")
    print(f"Documents without embeddings: {total_count - with_embedding_count}")

except Exception as e:
    print(f"An error occurred: {e}")

finally:
    # Close the connection
    if 'client' in locals():
        client.close()

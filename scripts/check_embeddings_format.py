from pymongo import MongoClient
from bson import ObjectId

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

try:
    # Connect to MongoDB
    client = MongoClient(connection_string)

    # Access the pr_analyzer database
    db = client["pr_analyzer"]

    # Access the pullrequests collection
    pullrequests = db["pullrequests"]

    # Query for a document with an embedding
    sample_doc = pullrequests.find_one({"embedding": {"$exists": True}})

    if sample_doc:
        print("Sample document found:")
        print(f"Document ID: {sample_doc['_id']}")
        embedding = sample_doc.get('embedding')
        if embedding:
            print(f"Embedding type: {type(embedding)}")
            print(f"Embedding length: {len(embedding)}")
        else:
            print("No embedding found in the document")
    else:
        print("No documents found with embeddings")

except Exception as e:
    print(f"An error occurred: {e}")

finally:
    # Close the connection
    if 'client' in locals():
        client.close()

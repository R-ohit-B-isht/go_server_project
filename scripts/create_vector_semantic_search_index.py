from pymongo import MongoClient, ASCENDING
from pymongo.operations import IndexModel

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

try:
    # Connect to MongoDB
    client = MongoClient(connection_string)

    # Access the pr_analyzer database
    db = client["pr_analyzer"]

    # Access the pullrequests collection
    pullrequests = db["pullrequests"]

    # Create the vectorSemanticSearch index
    index_model = IndexModel([("embedding", ASCENDING)], name="vectorSemanticSearch")
    pullrequests.create_indexes([index_model])

    print("vectorSemanticSearch index created successfully.")

    # Verify the index creation
    indexes = pullrequests.index_information()
    if "vectorSemanticSearch" in indexes:
        print("Verification: vectorSemanticSearch index exists.")
        print("Index details:")
        print(indexes["vectorSemanticSearch"])
    else:
        print("Verification failed: vectorSemanticSearch index was not created.")

except Exception as e:
    print(f"An error occurred: {e}")

finally:
    # Close the connection
    if 'client' in locals():
        client.close()
chan 

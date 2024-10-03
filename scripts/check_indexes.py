from pymongo import MongoClient
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

    # Get and print the indexes
    indexes = pullrequests.index_information()
    print("Indexes on the pullrequests collection:")
    for index_name, index_info in indexes.items():
        print(f"Index Name: {index_name}")
        print(f"Index Info: {index_info}")
        print("---")

    # Check for the vectorSemanticSearch index
    if "vectorSemanticSearch" in indexes:
        print("vectorSemanticSearch index exists.")
        print("vectorSemanticSearch index details:")
        print(indexes["vectorSemanticSearch"])
    else:
        print("vectorSemanticSearch index does not exist.")

except Exception as e:
    print(f"An error occurred: {e}")

finally:
    # Close the connection
    if 'client' in locals():
        client.close()

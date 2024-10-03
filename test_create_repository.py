import requests
import json
from pymongo import MongoClient
from bson import ObjectId

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

# Connect to MongoDB
client = MongoClient(connection_string)
db = client["pr_analyzer"]
repositories = db["repositories"]

# API endpoint
url = "http://localhost:8080/repositories"

# Test data
test_repo = {
    "name": "Test Repository",
    "url": "https://github.com/test/repo"
}

# Send POST request to create a new repository
response = requests.post(url, json=test_repo)

if response.status_code == 201:
    print("Repository created successfully")
    created_repo = response.json()

    # Get the ID of the created repository
    repo_id = created_repo.get("InsertedID")

    if repo_id:
        # Query the database to check the structure of the created repository
        repo_doc = repositories.find_one({"_id": ObjectId(repo_id)})

        if repo_doc:
            print("Repository document found in the database")

            # Check if pullRequests array exists and is empty
            if "pullRequests" in repo_doc and isinstance(repo_doc["pullRequests"], list) and len(repo_doc["pullRequests"]) == 0:
                print("pullRequests array is correctly initialized as an empty array")
            else:
                print("Error: pullRequests array is not correctly initialized")
        else:
            print("Error: Repository document not found in the database")
    else:
        print("Error: Failed to get the ID of the created repository")
else:
    print(f"Error creating repository: {response.status_code}")
    print(response.text)

# Clean up: Delete the test repository
if repo_id:
    delete_response = requests.delete(f"{url}/{repo_id}")
    if delete_response.status_code == 200:
        print("Test repository deleted successfully")
    else:
        print(f"Error deleting test repository: {delete_response.status_code}")

# Close the MongoDB connection
client.close()

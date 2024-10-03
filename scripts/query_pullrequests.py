import pymongo
from pymongo import MongoClient
from bson import ObjectId

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

# Connect to MongoDB
client = MongoClient(connection_string)

# Access the pr_analyzer database
db = client["pr_analyzer"]

# Access the pullrequests collection
pullrequests = db["pullrequests"]

# Query all documents in the pullrequests collection
cursor = pullrequests.find()

# Print the results
print("Pull Requests found:")
for document in cursor:
    print(f"ID: {document['_id']}")
    print(f"Repository ID: {document.get('repository_id', 'N/A')}")
    print(f"PR Number: {document.get('number', 'N/A')}")
    print(f"Title: {document.get('title', 'N/A')}")
    print("---")

# Query the repositories collection to get the pullRequests array
repositories = db["repositories"]
repo_cursor = repositories.find({}, {"_id": 1, "pullRequests": 1})

print("\nRepository Pull Request IDs:")
for repo in repo_cursor:
    print(f"Repository ID: {repo['_id']}")
    print("Pull Request IDs:")
    for pr_id in repo.get('pullRequests', []):
        print(f"  {pr_id}")
    print("---")

# Close the connection
client.close()

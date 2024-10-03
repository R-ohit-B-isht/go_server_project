import pymongo
from pymongo import MongoClient

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

# Connect to MongoDB
client = MongoClient(connection_string)

# Access the pr_analyzer database
db = client["pr_analyzer"]

# Access the repositories collection
repositories = db["repositories"]

# Query all documents in the repositories collection
cursor = repositories.find()

# Print the results
print("Repositories found:")
for document in cursor:
    print(document)

# Close the connection
client.close()

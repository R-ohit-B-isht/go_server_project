from pymongo import MongoClient
from bson.objectid import ObjectId

# MongoDB connection string
connection_string = "mongodb+srv://mentor:mentor@cluster0.hpj3khd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

try:
    # Connect to MongoDB
    client = MongoClient(connection_string)

    # Access the pr_analyzer database
    db = client["pr_analyzer"]

    # Access the repositories collection
    repositories = db["repositories"]

    # Query the most recently created repository
    repository = repositories.find_one(sort=[("_id", -1)])

    if repository:
        print("Repository found:")
        print(f"ID: {repository['_id']}")
        print(f"Name: {repository['name']}")
        print(f"URL: {repository['url']}")

        if 'pullRequests' in repository:
            print(f"pullRequests: {repository['pullRequests']}")
            if isinstance(repository['pullRequests'], list) and len(repository['pullRequests']) == 0:
                print("pullRequests is correctly initialized as an empty array.")
            else:
                print("Error: pullRequests is not an empty array.")
        else:
            print("Error: pullRequests field is missing.")
    else:
        print("No repository found.")

except Exception as e:
    print(f"An error occurred: {e}")

finally:
    # Close the connection
    if 'client' in locals():
        client.close()

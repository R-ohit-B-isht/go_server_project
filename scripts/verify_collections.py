import pymongo
from pymongo import MongoClient
from bson import ObjectId

# MongoDB connection string
connection_string = "mongodb://localhost:27017"

# Connect to MongoDB
client = MongoClient(connection_string)

# Access the pr_analyzer database
db = client["pr_analyzer"]

# Access the collections
pullrequests = db["pullrequests"]
repositories = db["repositories"]

# Count documents in pullrequests collection
pr_count = pullrequests.count_documents({})
print(f"Number of documents in pullrequests collection: {pr_count}")

# Count pullrequests in repositories collection
repo_pr_count = 0
repo_pr_ids = set()
for repo in repositories.find({}, {"pullRequests": 1}):
    repo_prs = repo.get("pullRequests", [])
    repo_pr_count += len(repo_prs)
    repo_pr_ids.update(repo_prs)

print(f"Number of pullRequests in repositories collection: {repo_pr_count}")

# Check for discrepancies
if pr_count != repo_pr_count:
    print(f"Discrepancy found: {pr_count - repo_pr_count} more documents in pullrequests collection")
else:
    print("No discrepancy found in count. Collections are consistent.")

# Check for matching ObjectIDs
pr_ids = set(pr["_id"] for pr in pullrequests.find({}, {"_id": 1}))
missing_in_repos = pr_ids - repo_pr_ids
missing_in_prs = repo_pr_ids - pr_ids

if missing_in_repos:
    print(f"Found {len(missing_in_repos)} pull requests in pullrequests collection not present in repositories:")
    for pr_id in missing_in_repos:
        print(f"  - {pr_id}")

if missing_in_prs:
    print(f"Found {len(missing_in_prs)} pull requests in repositories not present in pullrequests collection:")
    for pr_id in missing_in_prs:
        print(f"  - {pr_id}")

if not missing_in_repos and not missing_in_prs:
    print("All ObjectIDs match between collections.")

# Close the connection
client.close()

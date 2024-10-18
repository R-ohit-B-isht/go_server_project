as else #!/bin/bash

# Set environment variables
export PR_LOG_FILE_PATH="/home/ubuntu/go_server_project/pullrequest_collection.log"
export GITHUB_TOKEN="your_github_token_here"

# Set the API endpoint
API_ENDPOINT="http://localhost:8080/pullrequests/collect"

# Repository ID (replace with a valid ID from your database)
REPO_ID="66f1e55f5fb28a006018b775"

# Send the request
curl -X POST "${API_ENDPOINT}?id=${REPO_ID}" \
-H "Content-Type: application/json" \
-d '{
"startDate": "2023-01-01",
"endDate": "2023-12-31",
"dateFormat": "2006-01-02"
}'

# Check the log file
echo "Contents of ${PR_LOG_FILE_PATH}:"
tail -n 50 "${PR_LOG_FILE_PATH}"

# Use MongoDB CLI to check stored data (adjust as needed)
echo "Checking stored data in MongoDB:"
mongosh pr_analyzer --eval "db.pullrequests.find().limit(5).pretty()"

# Note: Make sure to replace 'your_github_token_here' with a valid GitHub token
# and adjust the REPO_ID to match a valid repository ID in your database.

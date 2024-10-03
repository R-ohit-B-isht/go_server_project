#!/bin/bash

# Set environment variables
export GITHUB_TOKEN="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"
export PR_LOG_FILE_PATH="/home/ubuntu/go_server_project/pullrequest_collection.log"
export PR_STATE="all"
export PR_PER_PAGE="50"

# Set the API endpoint
API_ENDPOINT="http://localhost:8080/pullrequests"

# Send the request
curl -X POST 'http://localhost:8080/pullrequests/collect?id=66f1e55f5fb28a006018b775' \
-H 'Content-Type: application/json' \
-H 'Origin: http://localhost:8080' \
--max-time 30 \
-d '{
"startDate": "2024-01-01",
"endDate": "2024-01-31",
"dateFormat": "2006-01-02"
}'

# Check the log file
echo "Contents of ${PR_LOG_FILE_PATH}:"
cat "${PR_LOG_FILE_PATH}"

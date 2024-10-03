This document provides instructions on how to use the new dynamic endpoint for collecting pull requests.

## Endpoint
`POST /pullrequests`

## Parameters
- **Query Parameter**: `id` - The repository ID.
- **JSON Body**:
  - `startDate`: The start date for fetching pull requests (format: `YYYY-MM-DD`).
  - `endDate`: The end date for fetching pull requests (format: `YYYY-MM-DD`).
  - `dateFormat`: The date format used (default: `2006-01-02`).

## Example Request
```bash
curl -X POST 'http://localhost:8080/pullrequests?id=66f1e55f5fb28a006018b775' \\
-H 'Content-Type: application/json' \\
-d '{
"startDate": "2024-01-01",
"endDate": "2024-01-31",
"dateFormat": "2006-01-02"
}'
```

## Expected Response
- Success: HTTP 200 OK with a message indicating successful data collection.
- Error: HTTP 400 Bad Request if parameters are missing or invalid.

## Error Handling
- Ensure the repository ID is valid and exists in the database.
- Check date formats and ensure they are correctly specified.

## Notes
- The endpoint uses the GitHub API to fetch pull request data.
- Ensure the server is running and connected to the database before making requests.

# Go Server Project Documentation

This document provides an overview of the models and routes in the Go server project, along with example `curl` requests for interacting with the API.

## Models

### Repository Model

- **Purpose**: Represents a repository with its details and associated pull requests.
- **Fields**:
  - `ID`: Unique identifier for the repository.
  - `Name`: Name of the repository.
  - `URL`: URL of the repository.
  - `LastProcessedDate`: Date when the repository was last processed.
  - `PullRequests`: Array of ObjectIDs referencing associated pull requests.

### PullRequest Model

- **Purpose**: Represents a pull request with its details and status.
- **Fields**:
  - `ID`: Unique identifier for the pull request.
  - `PRId`: Unique pull request ID.
  - `Repository`: ObjectID referencing the associated repository.
  - `Title`: Title of the pull request.
  - `Description`: Description of the pull request.
  - `Author`: Author of the pull request.
  - `CreatedAt`: Date when the pull request was created.
  - `LastUpdatedAt`: Date when the pull request was last updated.
  - `MergedAt`: Date when the pull request was merged.
  - `Status`: Status of the pull request (open, closed, merged).
  - `Labels`: Array of labels associated with the pull request.
  - `CustomTags`: Array of custom tags associated with the pull request.
  - `Complexity`: Complexity score of the pull request.
  - `TimeToMerge`: Time taken to merge the pull request.
  - `ConflictLikelihood`: Likelihood of conflicts in the pull request.
  - `SimilarityScore`: Similarity score of the pull request.
  - `Cluster`: ObjectID referencing the associated cluster.

### Cluster Model

- **Purpose**: Represents a cluster with its details and associated repository.
- **Fields**:
  - `ID`: Unique identifier for the cluster.
  - `Name`: Name of the cluster.
  - `Description`: Description of the cluster.
  - `Centroid`: Centroid data for the cluster.
  - `PRs`: Array of ObjectIDs referencing associated pull requests.
  - `ScoreAverage`: Average score of the cluster.
  - `Repository`: ObjectID referencing the associated repository.

## Routes

### Repository Routes

- **Create Repository**:
  - **Endpoint**: `POST /repositories`
  - **Example**:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"name":"TestRepo","url":"http://testrepo.com"}' http://localhost:8080/repositories
    ```

- **Get All Repositories**:
  - **Endpoint**: `GET /repositories`
  - **Example**:
    ```bash
    curl -X GET http://localhost:8080/repositories
    ```

- **Get Repository by ID**:
  - **Endpoint**: `GET /repositories/{id}`
  - **Example**:
    ```bash
    curl -X GET http://localhost:8080/repositories/{repository_id}
    ```

- **Update Repository**:
  - **Endpoint**: `PUT /repositories/{id}`
  - **Example**:
    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"name":"UpdatedRepo","url":"http://updatedrepo.com"}' http://localhost:8080/repositories/{repository_id}
    ```

- **Delete Repository**:
  - **Endpoint**: `DELETE /repositories/{id}`
  - **Example**:
    ```bash
    curl -X DELETE http://localhost:8080/repositories/{repository_id}
    ```

### PullRequest Routes

- **Create Pull Request**:
  - **Endpoint**: `POST /pullrequests`
  - **Example**:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"prId":"PR123","repository":"{repository_id}","title":"Test PR","author":"TestUser"}' http://localhost:8080/pullrequests
    ```

- **Get All Pull Requests**:
  - **Endpoint**: `GET /pullrequests`
  - **Example**:
    ```bash
    curl -X GET http://localhost:8080/pullrequests
    ```

- **Get Pull Request by ID**:
  - **Endpoint**: `GET /pullrequests/{id}`
  - **Example**:
    ```bash
    curl -X GET http://localhost:8080/pullrequests/{pullrequest_id}
    ```

- **Update Pull Request**:
  - **Endpoint**: `PUT /pullrequests/{id}`
  - **Example**:
    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"title":"Updated PR","author":"UpdatedUser"}' http://localhost:8080/pullrequests/{pullrequest_id}
    ```

- **Delete Pull Request**:
  - **Endpoint**: `DELETE /pullrequests/{id}`
  - **Example**:
    ```bash
    curl -X DELETE http://localhost:8080/pullrequests/{pullrequest_id}
    ```

- **Collect Pull Requests**:
  - **Endpoint**: `POST /pullrequests/collect`
  - **Description**: Collects pull requests from a specified repository within a given date range and stores them in the MongoDB database. Supports pagination to handle large datasets efficiently.
  - **Parameters**:
    - `id` (query): The repository ID for which pull requests are to be collected.
    - `startDate` (body): The start date for the date range (format: YYYY-MM-DD).
    - `endDate` (body): The end date for the date range (format: YYYY-MM-DD).
    - `pageSize` (body): The number of pull requests to process per page (default: 10).
    - `page` (body): The page number to process (default: 1).
  - **Response**:
    - `currentPage`: The current page number being processed.
    - `insertedCount`: The number of pull requests inserted into the database.
    - `message`: A success message indicating the operation's outcome.
    - `pageSize`: The number of pull requests processed per page.
    - `totalCount`: The total number of pull requests available for the specified date range.
  - **Example**:
    ```bash
    curl -X POST 'http://localhost:8080/pullrequests-collect?id=selectedRepo' \
    -H 'Content-Type: application/json' \
    -d '{
      "startDate": "2020-01-01",
      "endDate": "2021-12-31",
      "pageSize": 10,
      "page": 1
    }'

## Additional Information

- Replace `{repository_id}` and `{pullrequest_id}` with actual IDs from the database when using the example `curl` requests.
- Ensure the server is running on `localhost:8080` before making API requests.

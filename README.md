# Go Server with MongoDB Integration

This project is a Go server with MongoDB integration, designed to handle operations related to repositories, pull requests, clusters, contributors, analysis results, processing jobs, and configurations. The server uses a Bloom filter to optimize lookups for pull requests.

## Models

### Repository

- **Purpose**: Represents a code repository.
- **Fields**: `name`, `url`, `lastProcessedDate`, `pullRequests`
- **Example Data**:
  ```json
  {
    "name": "example-repo",
    "url": "https://github.com/example/repo",
    "lastProcessedDate": "2023-01-01T00:00:00Z",
    "pullRequests": ["60f004fdfe01cc4558a7d4a9", "60f004fdfe01cc4558a7d4b0"]
  }
  ```

### PullRequest

- **Purpose**: Represents a pull request in a repository.
- **Fields**: `prId`, `repository`, `title`, `description`, `author`, `createdAt`, `lastUpdatedAt`, `mergedAt`, `status`, `labels`, `customTags`, `complexity`, `timeToMerge`, `conflictLikelihood`, `similarityScore`, `cluster`
- **Example Data**:
  ```json
  {
    "prId": "12345",
    "repository": "60f004fdfe01cc4558a7d4a9",
    "title": "Add new feature",
    "description": "This PR adds a new feature.",
    "author": "Devin",
    "createdAt": "2023-01-01T00:00:00Z",
    "lastUpdatedAt": "2023-01-02T00:00:00Z",
    "status": "open"
  }
  ```

### Cluster

- **Purpose**: Represents a cluster of pull requests.
- **Fields**: `name`, `description`, `centroid`, `prs`, `scoreAverage`
- **Example Data**:
  ```json
  {
    "name": "Feature Cluster",
    "description": "Cluster of feature-related PRs",
    "centroid": {},
    "prs": ["60f004fdfe01cc4558a7d4a9"],
    "scoreAverage": 0.85
  }
  ```

### Contributor

- **Purpose**: Represents a contributor to the repository.
- **Fields**: `name`, `email`, `totalContributions`, `expertiseAreas`, `contributionsPerCluster`
- **Example Data**:
  ```json
  {
    "name": "Devin",
    "email": "devin@example.com",
    "totalContributions": 42,
    "expertiseAreas": ["Go", "MongoDB"],
    "contributionsPerCluster": [{"cluster": "60f004fdfe01cc4558a7d4a9", "count": 10}]
  }
  ```

### AnalysisResult

- **Purpose**: Represents the results of an analysis.
- **Fields**: `date`, `topClusters`, `trendAnalysis`, `contributorInsights`
- **Example Data**:
  ```json
  {
    "date": "2023-01-01T00:00:00Z",
    "topClusters": [{"cluster": "60f004fdfe01cc4558a7d4a9", "score": 0.9}],
    "trendAnalysis": {},
    "contributorInsights": {}
  }
  ```

### ProcessingJob

- **Purpose**: Represents a job for processing data.
- **Fields**: `startDate`, `endDate`, `status`, `progress`, `lastProcessedPR`, `filters`, `repositoriesProcessed`
- **Example Data**:
  ```json
  {
    "startDate": "2023-01-01T00:00:00Z",
    "status": "in-progress",
    "progress": 50,
    "lastProcessedPR": "60f004fdfe01cc4558a7d4a9",
    "filters": {},
    "repositoriesProcessed": ["60f004fdfe01cc4558a7d4a9"]
  }
  ```

### Configuration

- **Purpose**: Represents configuration settings for the application.
- **Fields**: `similarityThreshold`, `scoringSystem`, `filters`, `nlpSettings`, `bloomFilterSettings`
- **Example Data**:
  ```json
  {
    "similarityThreshold": 0.8,
    "scoringSystem": {},
    "filters": {},
    "nlpSettings": {},
    "bloomFilterSettings": {}
  }
  ```

## Routes

### Repository Routes

- **Create Repository**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"example-repo","url":"https://github.com/example/repo"}' http://localhost:8080/repositories
  ```

- **Get Repositories**:
  ```bash
  curl http://localhost:8080/repositories
  ```

### PullRequest Routes

- **Create Pull Request**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"prId":"12345","title":"Add new feature","author":"Devin"}' http://localhost:8080/pullrequests
  ```

- **Get Pull Requests**:
  ```bash
  curl http://localhost:8080/pullrequests
  ```

### Cluster Routes

- **Create Cluster**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"Feature Cluster","description":"Cluster of feature-related PRs"}' http://localhost:8080/clusters
  ```

- **Get Clusters**:
  ```bash
  curl http://localhost:8080/clusters
  ```

### Contributor Routes

- **Create Contributor**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"Devin","email":"devin@example.com"}' http://localhost:8080/contributors
  ```

- **Get Contributors**:
  ```bash
  curl http://localhost:8080/contributors
  ```

### AnalysisResult Routes

- **Create Analysis Result**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"date":"2023-01-01T00:00:00Z"}' http://localhost:8080/analysisresults
  ```

- **Get Analysis Results**:
  ```bash
  curl http://localhost:8080/analysisresults
  ```

### ProcessingJob Routes

- **Create Processing Job**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"startDate":"2023-01-01T00:00:00Z","status":"in-progress"}' http://localhost:8080/processingjobs
  ```

- **Get Processing Jobs**:
  ```bash
  curl http://localhost:8080/processingjobs
  ```

### Configuration Routes

- **Create Configuration**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"similarityThreshold":0.8}' http://localhost:8080/configurations
  ```

- **Get Configurations**:
  ```bash
  curl http://localhost:8080/configurations
  ```

## Notes

- Ensure MongoDB is running and accessible using the provided URI.
- The server listens on port 8080 by default.
- The Bloom filter is used to optimize pull request existence checks.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v39/github"
)

type PullRequest struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     string    `json:"state"`
	Author    string    `json:"author"`
}

func main() {
	owner := "MetaMask"
	repo := "metamask-extension"
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	client := github.NewClient(nil)
	ctx := context.Background()

	var allPRs []*github.PullRequest
	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			log.Fatalf("Error fetching pull requests: %v", err)
		}
		allPRs = append(allPRs, prs...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	var filteredPRs []PullRequest
	for _, pr := range allPRs {
		if pr.CreatedAt.After(startDate) && pr.CreatedAt.Before(endDate) {
			filteredPRs = append(filteredPRs, PullRequest{
				Number:    pr.GetNumber(),
				Title:     pr.GetTitle(),
				CreatedAt: pr.GetCreatedAt(),
				UpdatedAt: pr.GetUpdatedAt(),
				State:     pr.GetState(),
				Author:    pr.User.GetLogin(),
			})
		}
	}

	jsonData, err := json.MarshalIndent(filteredPRs, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile("github_prs.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("Fetched %d pull requests for the specified date range.\n", len(filteredPRs))
	fmt.Println("Data saved to github_prs.json")
}

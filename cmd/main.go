package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type MilestoneData struct {
	Milestones []Milestone `yaml:"milestones"`
}

type Milestone struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	DueOn       string  `yaml:"due_on"`
	Issues      []Issue `yaml:"issues"`
}

type Issue struct {
	Title  string   `yaml:"title"`
	Body   string   `yaml:"body"`
	Labels []string `yaml:"labels"`
}

func createGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func loadMilestoneData(filename string) (*MilestoneData, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var data MilestoneData
	if err := yaml.Unmarshal(fileContent, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &data, nil
}

func createMilestone(ctx context.Context, client *github.Client, owner, repo string, ms Milestone) (*github.Milestone, error) {
	dueDate, err := time.Parse("2006-01-02", ms.DueOn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse due date %s: %w", ms.DueOn, err)
	}

	milestone, _, err := client.Issues.CreateMilestone(ctx, owner, repo, &github.Milestone{
		Title:       github.String(ms.Title),
		Description: github.String(ms.Description),
		DueOn:       &github.Timestamp{Time: dueDate},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create milestone %s: %w", ms.Title, err)
	}

	return milestone, nil
}

func createIssues(ctx context.Context, client *github.Client, owner, repo string, issues []Issue, milestoneNumber *int) error {
	for _, issue := range issues {
		_, _, err := client.Issues.Create(ctx, owner, repo, &github.IssueRequest{
			Title:     github.String(issue.Title),
			Body:      github.String(issue.Body),
			Milestone: milestoneNumber,
			Labels:    &issue.Labels,
		})
		if err != nil {
			return fmt.Errorf("failed to create issue %s: %w", issue.Title, err)
		}
	}
	return nil
}

func main() {
	ctx := context.Background()

	// Configuration
	const (
		tokenPlaceholder = "YOUR_GITHUB_TOKEN_HERE"
		yamlFile         = "milestone.yaml"
		owner            = "ORGANIZATION_NAME"
		repo             = "REPOSITORY_NAME"
	)

	client := createGitHubClient(ctx, tokenPlaceholder)

	data, err := loadMilestoneData(yamlFile)
	if err != nil {
		fmt.Printf("Error loading milestone data: %v\n", err)
		return
	}

	for _, ms := range data.Milestones {
		milestone, err := createMilestone(ctx, client, owner, repo, ms)
		if err != nil {
			fmt.Printf("Error creating milestone: %v\n", err)
			return
		}

		if err := createIssues(ctx, client, owner, repo, ms.Issues, milestone.Number); err != nil {
			fmt.Printf("Error creating issues: %v\n", err)
			return
		}
	}
}

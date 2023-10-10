package api

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v55/github"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type PRWithStatus struct {
	PR     *github.PullRequest
	Status string
}

const pageSizeMax = 100

// Singleton
var githubClient *github.Client

func getAccessToken() (*oauth2.Token, error) {
	godotenv.Load()
	tokenString, isSet := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !isSet {
		return nil, fmt.Errorf("error: Could not find GITHUB_ACCESS_TOKEN. Make sure token is set in env or add it to a .env file in your project root")
	}

	token := oauth2.Token{AccessToken: tokenString}

	return &token, nil
}

func getClientInstance() (*github.Client, error) {
	if githubClient == nil {
		tokenString, err := getAccessToken()
		if err != nil {
			return nil, err
		}
		tokenSource := oauth2.StaticTokenSource(tokenString)
		httpClient := oauth2.NewClient(context.Background(), tokenSource)

		githubClient = github.NewClient(httpClient)
	}

	return githubClient, nil
}

func parseRepoPath(repoPath string) (string, string, error) {
	var owner, repo string

	path := strings.Split(repoPath, "/")
	if len(path) != 2 {
		return "", "", fmt.Errorf("error: invalid repo path '%s'", repoPath)
	}

	owner = path[0]
	repo = path[1]

	return owner, repo, nil
}

func GetOwnedRepos(username string) ([]*github.Repository, error) {
	client, err := getClientInstance()
	if err != nil {
		return nil, err
	}

	repos, _, err := client.Repositories.List(context.Background(), username, nil)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func GetFollowedRepos(username string) ([]*github.Repository, error) {
	client, err := getClientInstance()
	if err != nil {
		return nil, err
	}

	repos, _, err := client.Activity.ListWatched(context.Background(), username, nil)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func ListRepoWorkflows(repoPath string) (*github.Workflows, error) {
	client, err := getClientInstance()
	if err != nil {
		return nil, err
	}

	owner, repo, err := parseRepoPath(repoPath)
	if err != nil {
		return nil, err
	}

	workflows, _, err := client.Actions.ListWorkflows(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

func ListPRsByRepo(repoPath string) ([]*github.PullRequest, error) {
	client, err := getClientInstance()
	if err != nil {
		return nil, err
	}

	owner, repo, err := parseRepoPath(repoPath)
	if err != nil {
		return nil, err
	}

	options := github.PullRequestListOptions{State: "open", Sort: "created", Direction: "desc"}
	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, &options)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func GetPRStatus(repoPath string, pr *github.PullRequest) (string, error) {
	if pr == nil {
		return "", fmt.Errorf("error: invalid pull request")
	}

	client, err := getClientInstance()
	if err != nil {
		return "", err
	}

	owner, repo, err := parseRepoPath(repoPath)
	if err != nil {
		return "", err
	}

	statuses, _, err := client.Repositories.GetCombinedStatus(context.Background(), owner, repo, *pr.Head.SHA, nil)
	if err != nil {
		return "", err
	}

	return statuses.GetState(), nil
}

func ListPRsByRepoWithStatus(repoPath string) ([]*PRWithStatus, error) {
	prs, err := ListPRsByRepo(repoPath)
	if err != nil {
		return nil, err
	}

	var prsWithStatus []*PRWithStatus
	for _, pr := range prs {
		status, err := GetPRStatus(repoPath, pr)
		if err != nil {
			return nil, err
		}
		prsWithStatus = append(prsWithStatus, &PRWithStatus{PR: pr, Status: status})
	}

	return prsWithStatus, nil
}

func ListPRsByAuthor(author string, size int) ([]*github.Issue, error) {
	client, err := getClientInstance()
	if err != nil {
		return nil, err
	}

	page := 1
	pageSize := pageSizeMax
	var prs []*github.Issue

	for size > 0 {
		if size < pageSize {
			pageSize = size
		}

		options := &github.SearchOptions{Sort: "created", Order: "desc", ListOptions: github.ListOptions{Page: page, PerPage: pageSize}}
		res, _, err := client.Search.Issues(context.Background(), fmt.Sprintf("is:pr author:%s status:failure", author), options)
		if err != nil {
			return nil, err
		}

		prs = append(prs, res.Issues...)

		size -= pageSize
		page += 1
	}

	return prs, nil
}

package api

import (
	"os"
	"testing"

	"github.com/google/go-github/v55/github"
	"github.com/joho/godotenv"
)

const (
	expectedErrorGotNil    = "expected error, got nil"
	expectedDifferentError = "expected error '%s', got '%s'"
	expectedNoError        = "expected no error, but got '%s'"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	code := m.Run()
	os.Exit(code)
}

func TestGetAccessTokenWithNoAccessTokenSet(t *testing.T) {
	loadErr := godotenv.Load()
	if loadErr == nil {
		os.Rename(".env", ".env.bak")
	}

	token, isSet := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if isSet {
		os.Unsetenv("GITHUB_ACCESS_TOKEN")
	}

	_, err := getAccessToken()

	expectedError := "could not find GITHUB_ACCESS_TOKEN. Make sure token is set in env or add it to a .env file in your project root"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {

		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}

	t.Cleanup(func() {
		if loadErr == nil {
			os.Rename(".env.bak", ".env")
		}
		if isSet {
			os.Setenv("GITHUB_ACCESS_TOKEN", token)
		}
	})
}

func TestGetAccessTokenWithAccessTokenSet(t *testing.T) {
	token, err := getAccessToken()

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if token == nil {
		t.Fatal("expected a non-nil token, but got nil")
	}

	tokenString := token.AccessToken
	expectedTokenString := os.Getenv("GITHUB_ACCESS_TOKEN")

	if tokenString != expectedTokenString {
		t.Errorf("expected token does not match actual token")
	}
}

func TestGetClientInstance(t *testing.T) {
	client, err := getClientInstance()

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if client == nil {
		t.Fatal("expected a non-nil GitHub client, but got nil")
	}

	baseURL := client.BaseURL.String()
	expectedBaseURL := "https://api.github.com/"

	if baseURL != expectedBaseURL {
		t.Errorf("expected client base URL to be %s, but got %s", expectedBaseURL, baseURL)
	}
}

func TestParseRepoPathWithEmptyPath(t *testing.T) {
	_, _, err := parseRepoPath("")

	expectedError := "invalid repo path ''"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestParseRepoPathWithInvalidPathShort(t *testing.T) {
	_, _, err := parseRepoPath("notavalidpath")

	expectedError := "invalid repo path 'notavalidpath'"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestParseRepoPathWithInvalidPathLong(t *testing.T) {
	_, _, err := parseRepoPath("not/a/valid/path")

	expectedError := "invalid repo path 'not/a/valid/path'"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestParseRepoWithValidPath(t *testing.T) {
	owner, repo, err := parseRepoPath("owner/repo")

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if owner != "owner" {
		t.Errorf("expected owner to be 'owner', but got '%s'", owner)
	}

	if repo != "repo" {
		t.Errorf("expected repo to be 'repo', but got '%s'", repo)
	}
}

func TestGetOwnedReposWithInvalidUsername(t *testing.T) {
	_, err := GetOwnedRepos("gidhjfgu90w45u", 30)

	expectedError := "could not retrieve repositories for user 'gidhjfgu90w45u', make sure the username is valid and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetOwnedReposWithValidUsername(t *testing.T) {
	expectedName := "carolinafsilva"
	repos, err := GetOwnedRepos(expectedName, 30)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if repos == nil {
		t.Fatal("expected a non-nil list of repos, but got nil")
	}

	if len(repos) == 0 {
		t.Fatal("expected a non-empty list of repos, but got an empty list")
	}

	for _, repo := range repos {
		if repo == nil {
			t.Fatal("expected a non-nil repo, but got nil")
		}

		name := repo.GetOwner().GetLogin()

		if name != expectedName {
			t.Errorf("expected owner to be '%s', but got '%s'", expectedName, name)
		}
	}
}

func TestGetFollowedReposWithInvalidUsername(t *testing.T) {
	_, err := GetFollowedRepos("gidhjfgu90w45u", 30)

	expectedError := "could not retrieve followed repositories for user 'gidhjfgu90w45u', make sure the username is valid and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetFollowedReposWithValidUsername(t *testing.T) {
	expectedName := "carolinafsilva"
	repos, err := GetFollowedRepos(expectedName, 30)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if repos == nil {
		t.Fatal("expected a non-nil list of repos, but got nil")
	}

	if len(repos) == 0 {
		t.Fatal("expected a non-empty list of repos, but got an empty list")
	}

	for _, repo := range repos {
		if repo == nil {
			t.Fatal("expected a non-nil repo, but got nil")
		}
	}
}

func TestListRepoWorkflowsWithInvalidRepoPath(t *testing.T) {
	_, err := ListRepoWorkflows("notavalidpath")

	expectedError := "invalid repo path 'notavalidpath'"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListRepoWorkflowsWithInvalidOwner(t *testing.T) {
	_, err := ListRepoWorkflows("gidhjfgu90w45u/repo")

	expectedError := "could not retrieve workflows for repo 'gidhjfgu90w45u/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListRepoWorkflowsWithInvalidRepo(t *testing.T) {
	_, err := ListRepoWorkflows("carolinafsilva/repo")

	expectedError := "could not retrieve workflows for repo 'carolinafsilva/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListRepoWorkflowsWithValidRepoPath(t *testing.T) {
	repoPath := "aleph-two/flowcar.pt"
	workflows, err := ListRepoWorkflows(repoPath)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if workflows == nil {
		t.Fatal("expected a non-nil list of workflows, but got nil")
	}

	if len(workflows.Workflows) == 0 {
		t.Fatal("expected a non-empty list of workflows, but got an empty list")
	}

	for _, workflow := range workflows.Workflows {
		if workflow == nil {
			t.Fatal("expected a non-nil workflow, but got nil")
		}
	}
}

func TestListPRsByRepoWithInvalidRepoPath(t *testing.T) {
	_, err := ListPRsByRepo("notavalidpath", 30)

	expectedError := "invalid repo path 'notavalidpath'"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByRepoWithInvalidOwner(t *testing.T) {
	_, err := ListPRsByRepo("gidhjfgu90w45u/repo", 30)

	expectedError := "could not retrieve pull requests for repo 'gidhjfgu90w45u/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByRepoWithInvalidRepo(t *testing.T) {
	_, err := ListPRsByRepo("carolinafsilva/repo", 30)

	expectedError := "could not retrieve pull requests for repo 'carolinafsilva/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByRepoWithValidRepoPath(t *testing.T) {
	repoPath := "aleph-two/flowcar.pt"
	prs, err := ListPRsByRepo(repoPath, 30)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if prs == nil {
		t.Fatal("expected a non-nil list of PRs, but got nil")
	}

	if len(prs) == 0 {
		t.Fatal("expected a non-empty list of PRs, but got an empty list")
	}

	for _, pr := range prs {
		if pr == nil {
			t.Fatal("expected a non-nil PR, but got nil")
		}
	}
}

func TestGetPRStatusWithInvalidRepoPath(t *testing.T) {
	pr := &github.PullRequest{}
	pr.Head = &github.PullRequestBranch{}
	pr.Head.SHA = github.String("1edbbf3b63d57d8f4f22e1c4617aa2e2ca4c7d96")

	_, err := GetPRStatus("notavalidpath", pr)

	expectedError := "invalid repo path 'notavalidpath'"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetPRStatusWithInvalidOwner(t *testing.T) {
	pr := &github.PullRequest{}
	pr.Head = &github.PullRequestBranch{}
	pr.Head.SHA = github.String("1edbbf3b63d57d8f4f22e1c4617aa2e2ca4c7d96")

	_, err := GetPRStatus("gidhjfgu90w45u/repo", pr)

	expectedError := "could not retrieve status for pull request in 'gidhjfgu90w45u/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetPRStatusWithInvalidRepo(t *testing.T) {
	pr := &github.PullRequest{}
	pr.Head = &github.PullRequestBranch{}
	pr.Head.SHA = github.String("1edbbf3b63d57d8f4f22e1c4617aa2e2ca4c7d96")

	_, err := GetPRStatus("carolinafsilva/repo", pr)

	expectedError := "could not retrieve status for pull request in 'carolinafsilva/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetPRStatusWithNillPR(t *testing.T) {
	_, err := GetPRStatus("aleph-two/flowcar.pt", nil)

	expectedError := "invalid pull request"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetPRStatusWithInvalidPR(t *testing.T) {
	pr := &github.PullRequest{}
	pr.Head = &github.PullRequestBranch{}
	pr.Head.SHA = github.String("82758932759379857349859835473498")

	_, err := GetPRStatus("aleph-two/flowcar.pt", pr)

	expectedError := "could not retrieve status for pull request in 'aleph-two/flowcar.pt', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestGetPRStatusWithValidPR(t *testing.T) {
	pr := &github.PullRequest{}
	pr.Head = &github.PullRequestBranch{}
	pr.Head.SHA = github.String("1edbbf3b63d57d8f4f22e1c4617aa2e2ca4c7d96")

	status, err := GetPRStatus("aleph-two/flowcar.pt", pr)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if status == "" {
		t.Fatal("expected a non-empty status, but got an empty string")
	}
}

func TestListPRsByRepoWithStatusWithInvalidPath(t *testing.T) {
	_, err := ListPRsByRepoWithStatus("notavalidpath", 30)

	expectedError := "invalid repo path 'notavalidpath'"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByRepoWithStatusWithInvalidOwner(t *testing.T) {
	_, err := ListPRsByRepoWithStatus("gidhjfgu90w45u/repo", 30)

	expectedError := "could not retrieve pull requests for repo 'gidhjfgu90w45u/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByRepoWithStatusWithInvalidRepo(t *testing.T) {
	_, err := ListPRsByRepoWithStatus("carolinafsilva/repo", 30)

	expectedError := "could not retrieve pull requests for repo 'carolinafsilva/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByRepoWithStatusWithValidRepoPath(t *testing.T) {
	repoPath := "aleph-two/flowcar.pt"
	prsWithStatus, err := ListPRsByRepoWithStatus(repoPath, 30)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if prsWithStatus == nil {
		t.Fatal("expected a non-nil list of PRs, but got nil")
	}

	if len(prsWithStatus) == 0 {
		t.Fatal("expected a non-empty list of PRs, but got an empty list")
	}

	for _, pr := range prsWithStatus {
		if pr == nil {
			t.Fatal("expected a non-nil PR, but got nil")
		}

		if pr.Status == "" {
			t.Fatal("expected a non-empty status, but got an empty string")
		}
	}
}

func TestListPRsByAuthorWithInvalidAuthor(t *testing.T) {
	_, err := ListPRsByAuthor("gidhjfgu90w45u", 30)

	expectedError := "could not retrieve pull requests for author 'gidhjfgu90w45u', make sure the username is valid and GITHUB_ACCESS_TOKEN is set and valid"
	if err == nil {
		t.Error(expectedErrorGotNil)
	} else if err.Error() != expectedError {
		t.Errorf(expectedDifferentError, expectedError, err.Error())
	}
}

func TestListPRsByAuthorWithValidAuthor(t *testing.T) {
	prs, err := ListPRsByAuthor("willnorris", 30)

	if err != nil {
		t.Errorf(expectedNoError, err.Error())
	}

	if prs == nil {
		t.Fatal("expected a non-nil list of PRs, but got nil")
	}

	if len(prs) == 0 {
		t.Fatal("expected a non-empty list of PRs, but got an empty list")
	}

	for _, pr := range prs {
		if pr == nil {
			t.Fatal("expected a non-nil PR, but got nil")
		}
	}
}

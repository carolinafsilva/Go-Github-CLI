package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/carolinafsilva/go-github-cli/api"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	os.Exit(m.Run())
}

func TestPrCmdPrintHelpMenu(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	args := []string{"--help", "-h", "help", "", "invalidcmd"}

	expectedOutput := "The pr command in GG is designed to retrieve essential pull request information from GitHub.\n\tYou can use this command to filter and display pull requests based on different criteria such as the author or repository\n\nUsage:\n  gg pr [command]\n\nAvailable Commands:\n  author      Get Pull Request information by author\n  repo        Get Pull Request information by repository\n\nFlags:\n  -h, --help   help for pr\n\nUse \"gg pr [command] --help\" for more information about a command.\n"

	for _, arg := range args {
		cmd.SetArgs([]string{"pr", arg})
		err := cmd.Execute()
		if err != nil {
			t.Errorf(expectedNoError, err)
		}

		if output.String() != expectedOutput {
			t.Errorf(expectedDifferentError, expectedOutput, output.String())
		}

		output.Reset()
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestPrCmdWithInvalidShorthandFlag(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"pr", "-i"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "unknown shorthand flag: 'i' in -i"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrAuthorCmdWithNoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"pr", "author"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 0"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrAuthorCmdWithTwoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"pr", "author", "carolina", "silva"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 2"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrAuthorCmdWithInvalidUsername(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"pr", "author", "gidhjfgu90w45u"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := "could not retrieve pull requests for author 'gidhjfgu90w45u', make sure the username is valid and GITHUB_ACCESS_TOKEN is set and valid\n"
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestPrAuthorCmdWithValidUsername(t *testing.T) {
	cmd := rootCmd

	author := "carolinafsilva"
	size := 30

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"pr", "author", author})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	pr, _ := api.ListPRsByAuthor(author, size)

	var expectedMsg string
	for i, pr := range pr {
		expectedMsg += fmt.Sprintf("%3d. %s %s\n", i+1, *pr.CreatedAt, *pr.Title)
	}
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestPrAuthorCmdWithValidUsernameAndSizeFlagWithNoArg(t *testing.T) {
	cmd := rootCmd

	author := "carolinafsilva"

	cmd.SetArgs([]string{"pr", "author", author, "-S"})

	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "flag needs an argument: 'S' in -S"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrAuthorCmdWithValidUsernameAndSizeFlagWithInvalidSize(t *testing.T) {
	cmd := rootCmd

	author := "carolinafsilva"

	cmd.SetArgs([]string{"pr", "author", author, "-S", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "invalid argument \"invalid\" for \"-S, --size\" flag: strconv.ParseInt: parsing \"invalid\": invalid syntax"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrAuthorCmdWithValidUsernameAndSizeFlagWithValidSize(t *testing.T) {
	cmd := rootCmd

	author := "carolinafsilva"
	size := 30

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"pr", "author", author, "-S", "30"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	pr, _ := api.ListPRsByAuthor(author, size)

	var expectedMsg string
	for i, pr := range pr {
		expectedMsg += fmt.Sprintf("%3d. %s %s\n", i+1, *pr.CreatedAt, *pr.Title)
	}
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestPrRepoCmdWithNoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"pr", "repo"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 0"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrRepoCmdWithTwoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"pr", "repo", "go", "github"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 2"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestPrRepoCmdWithInvalidRepoPath(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"pr", "repo", "carolinafsilva/repo"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := "could not retrieve pull requests for repo 'carolinafsilva/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid\n"
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestPrRepoCmdWithValidRepoPath(t *testing.T) {
	cmd := rootCmd

	repoPath := "carolinafsilva/go-github-cli"

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"pr", "repo", repoPath})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	pr, _ := api.ListPRsByRepo(repoPath, 30)

	var expectedMsg string
	for i, pr := range pr {
		expectedMsg += fmt.Sprintf("%3d. %s %s\n", i+1, *pr.CreatedAt, *pr.Title)
	}
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestPrRepoCmdWithValidRepoPathAndStatusFlag(t *testing.T) {
	cmd := rootCmd

	repoPath := "carolinafsilva/go-github-cli"

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"pr", "repo", repoPath, "--status"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	prs, _ := api.ListPRsByRepoWithStatus(repoPath, 30)

	var expectedMsg string
	for i, pr := range prs {
		expectedMsg += fmt.Sprintf("%3d. %s  %s  %s\n", i+1, *pr.PR.CreatedAt, pr.Status, *pr.PR.Title)
	}
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

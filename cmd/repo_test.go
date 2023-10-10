package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/carolinafsilva/go-github-cli/api"
)

const (
	expectedNoError        = "expected no error, got %v"
	expectedDifferentError = "expected %v, got %v"
	expectedErrorGotNil    = "expected error, got nil"
)

func TestRepoCmdPrintHelpMenu(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	args := []string{"--help", "-h", "help", "", "invalidcmd"}

	expectedOutput := "The repo command in GG allows you to interact with GitHub repositories. This command provides subcommands for listing repositories and their workflows.\n\nUsage:\n  gg repo [command]\n\nAvailable Commands:\n  list        List a user's repositories\n  workflow    List a repository's workflows\n\nFlags:\n  -h, --help   help for repo\n\nUse \"gg repo [command] --help\" for more information about a command.\n"

	for _, arg := range args {
		cmd.SetArgs([]string{"repo", arg})
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

func TestRepoCmdWithInvalidShorthandFlag(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"repo", "-i"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "unknown shorthand flag: 'i' in -i"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestRepoListCmdWithNoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"repo", "list"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 0"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestRepoListCmdWithTwoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"repo", "list", "foo", "bar"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 2"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestRepoListCmdWithInvalidUsername(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"repo", "list", "gidhjfgu90w45u"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := "could not retrieve repositories for user 'gidhjfgu90w45u', make sure the username is valid and GITHUB_ACCESS_TOKEN is set and valid\n"
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestRepoListCmdWithValidUsername(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"repo", "list", "carolinafsilva"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := fmt.Sprintln("Owned Repositories:")
	repos, _ := api.GetOwnedRepos("carolinafsilva", 30)
	for _, repo := range repos {
		expectedMsg += fmt.Sprintln(*repo.Name)
	}

	expectedMsg += fmt.Sprintln("Followed Repositories:")
	repos, _ = api.GetFollowedRepos("carolinafsilva", 30)
	for _, repo := range repos {
		expectedMsg += fmt.Sprintln(*repo.Name)
	}

	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestRepoListCmdWithValidUsernameAndOwnedFlag(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"repo", "list", "carolinafsilva", "--owned"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := fmt.Sprintln("Owned Repositories:")
	repos, _ := api.GetOwnedRepos("carolinafsilva", 30)
	for _, repo := range repos {
		expectedMsg += fmt.Sprintln(*repo.Name)
	}

	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
		repoListCmd.ResetFlags()
		repoListCmd.Flags().BoolVar(&owned, "owned", false, "List only owned repos")
		repoListCmd.Flags().BoolVar(&followed, "followed", false, "List only followed repos")
		repoListCmd.MarkFlagsMutuallyExclusive("owned", "followed")
	})
}

func TestRepoListCmdWithValidUsernameAndFollowedFlag(t *testing.T) {
	cmd := rootCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"repo", "list", "carolinafsilva", "--followed"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := fmt.Sprintln("Followed Repositories:")
	repos, _ := api.GetFollowedRepos("carolinafsilva", 30)
	for _, repo := range repos {
		expectedMsg += fmt.Sprintln(*repo.Name)
	}

	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
		repoListCmd.ResetFlags()
		repoListCmd.Flags().BoolVar(&owned, "owned", false, "List only owned repos")
		repoListCmd.Flags().BoolVar(&followed, "followed", false, "List only followed repos")
		repoListCmd.MarkFlagsMutuallyExclusive("owned", "followed")
	})
}

func TestRepoWorkflowCmdWithNoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"repo", "workflow"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 0"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestRepoWorkflowCmdWithTwoArgs(t *testing.T) {
	cmd := rootCmd

	cmd.SetArgs([]string{"repo", "workflow", "foo", "bar"})
	err := cmd.Execute()
	if err == nil {
		t.Errorf(expectedErrorGotNil)
	}

	expectedErr := "accepts 1 arg(s), received 2"
	if err.Error() != expectedErr {
		t.Errorf(expectedDifferentError, expectedErr, err.Error())
	}
}

func TestRepoWorkflowCmdWithInvalidRepoPath(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"repo", "workflow", "carolinafsilva/repo"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	expectedMsg := "could not retrieve workflows for repo 'carolinafsilva/repo', make sure the repository exists and GITHUB_ACCESS_TOKEN is set and valid\n"
	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

func TestRepoWorkflowCmdWithValidRepoPath(t *testing.T) {
	cmd := rootCmd

	var output bytes.Buffer
	cmd.SetOut(&output)

	cmd.SetArgs([]string{"repo", "workflow", "aleph-two/flowcar.pt"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf(expectedNoError, err)
	}

	workflows, _ := api.ListRepoWorkflows("aleph-two/flowcar.pt")

	var expectedMsg string
	for _, workflow := range workflows.Workflows {
		expectedMsg += fmt.Sprintf("%s\n", workflow.GetName())
	}

	if output.String() != expectedMsg {
		t.Errorf(expectedDifferentError, expectedMsg, output.String())
	}

	t.Cleanup(func() {
		cmd.SetOut(nil)
	})
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var githubUser string

var rootCmd = &cobra.Command{
	Use:   "gg <command> [subcommand] [flags]",
	Short: "gg is a command-line tool for interacting with GitHub's Pull Requests and Repositories",
	Long:  `gg is a versatile command-line tool for interacting with GitHub. It provides subcommands to access information about GitHub Pull Requests and Repositories`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

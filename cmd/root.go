package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	githubUser string
	fg         = color.New()
	magenta    = color.New(color.FgMagenta)
)

var rootCmd = &cobra.Command{
	Use:   "gg <command> [subcommand] [flags]",
	Short: "gg is a command-line tool for interacting with GitHub's Pull Requests and Repositories",
	Long:  `gg is a versatile command-line tool for interacting with GitHub. It provides subcommands to access information about GitHub Pull Requests and Repositories`,
}

func Execute() {
	color.NoColor = false
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

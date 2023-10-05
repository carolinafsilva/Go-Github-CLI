package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var githubUser string

var rootCmd = &cobra.Command{
	Use:   "gg <command> [subcommand] [flags]",
	Short: "TODO change short description",
	Long:  `TODO change long description`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

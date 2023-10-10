package cmd

import (
	"fmt"

	"github.com/carolinafsilva/go-github-cli/api"
	"github.com/spf13/cobra"
)

var (
	size   int
	status bool
)

var prCmd = &cobra.Command{
	Use:   "pr <command> [flags]",
	Short: "TODO - Add short description of the command here",
	Long:  `TODO - Add long description of the command here`,
}

var prAuthorCmd = &cobra.Command{
	Use:   "author <username> [flags]",
	Short: "TODO - Add short description of the command here",
	Long:  `TODO - Add long description of the command here`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		author := args[0]

		prs, err := api.ListPRsByAuthor(author, size)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, pr := range prs {
			fmt.Println(*pr.Title, *pr.CreatedAt)
		}
	},
}

var prRepoCmd = &cobra.Command{
	Use:   "repo <owner/repo> [flags]",
	Short: "TODO - Add short description of the command here",
	Long:  `TODO - Add long description of the command here`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]

		if status {
			prs, err := api.ListPRsByRepoWithStatus(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, pr := range prs {
				fmt.Printf("%s: %s\n", *pr.PR.Title, pr.Status)
			}
		} else {
			prs, err := api.ListPRsByRepo(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, pr := range prs {
				fmt.Println(*pr.Title, *pr.CreatedAt)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(prAuthorCmd)
	prCmd.AddCommand(prRepoCmd)

	prAuthorCmd.Flags().IntVarP(&size, "size", "S", 30, "Number of results to return")
	prRepoCmd.Flags().BoolVarP(&status, "status", "s", false, "Show the state of the PRs in the Workflow")
}

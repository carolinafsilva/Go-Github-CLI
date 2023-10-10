package cmd

import (
	"fmt"

	"github.com/carolinafsilva/go-github-cli/api"
	"github.com/spf13/cobra"
)

var (
	owned    bool
	followed bool
)

var repoCmd = &cobra.Command{
	Use:   "repo [command]",
	Short: "TODO - Add short description",
	Long:  `TODO - Add long description`,
}

var repoListCmd = &cobra.Command{
	Use:   "list <username>",
	Short: "TODO - Add short description",
	Long:  `TODO - Add long description`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		githubUser = args[0]

		if !followed {
			owned, err := api.GetOwnedRepos(githubUser)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("\n***OWNED REPOS***")
			for _, repo := range owned {
				fmt.Println(*repo.Name)
			}
		}

		if !owned {
			followed, err := api.GetFollowedRepos(githubUser)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("\n***FOLLOWED REPOS***")
			for _, repo := range followed {
				fmt.Println(*repo.Name)
			}
		}
	},
}

var repoWorkflowCmd = &cobra.Command{
	Use:   "workflow <owner/repo>",
	Args:  cobra.ExactArgs(1),
	Short: "TODO - Add short description",
	Long:  `TODO - Add long description`,
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]

		workflows, err := api.ListRepoWorkflows(repoPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(workflows.Workflows) == 0 {
			fmt.Println("The repository does not have workflows.")
			return
		}

		for _, workflow := range workflows.Workflows {
			fmt.Printf("%s (#%d)\n", workflow.GetName(), workflow.GetID())
		}

	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoWorkflowCmd)

	repoListCmd.Flags().BoolVar(&owned, "owned", false, "List only owned repos")
	repoListCmd.Flags().BoolVar(&followed, "followed", false, "List only followed repos")
	repoListCmd.MarkFlagsMutuallyExclusive("owned", "followed")
}

package cmd

import (
	"github.com/carolinafsilva/go-github-cli/api"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	owned    bool
	followed bool
)

var repoCmd = &cobra.Command{
	Use:   "repo [command]",
	Short: "Get information about Github Repositories",
	Long:  `The repo command in GG allows you to interact with GitHub repositories. This command provides subcommands for listing repositories and their workflows.`,
}

var repoListCmd = &cobra.Command{
	Use:   "list <username>",
	Short: "List a user's repositories",
	Long: `The list subcommand within the repo command lists GitHub repositories based on specified criteria. You can use this command to retrieve repositories owned or followed by a particular user. Additionally, you can filter the results to display only owned or followed repositories by using the respective flags.

--followed: List only followed repositories.
--owned: List only owned repositories.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		githubUser = args[0]

		magenta := color.New(color.FgMagenta).Add(color.Underline)

		if !followed {
			repos, err := api.GetOwnedRepos(githubUser, size)
			if err != nil {
				cmd.Println(err)
				return
			}

			magenta.Fprintln(cmd.OutOrStdout(), "Owned Repositories:")
			for _, repo := range repos {
				cmd.Println(*repo.Name)
			}
		}

		if !owned {
			repos, err := api.GetFollowedRepos(githubUser, size)
			if err != nil {
				cmd.Println(err)
				return
			}

			magenta.Fprintln(cmd.OutOrStdout(), "Followed Repositories:")
			for _, repo := range repos {
				cmd.Println(*repo.Name)
			}
		}
	},
}

var repoWorkflowCmd = &cobra.Command{
	Use:   "workflow <owner/repo>",
	Args:  cobra.ExactArgs(1),
	Short: "List a repository's workflows",
	Long:  `The workflow subcommand within the repo command allows you to access and view the workflows associated with a specific GitHub repository. You can specify the repository using the <owner/repo> parameter to retrieve information about its workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]

		workflows, err := api.ListRepoWorkflows(repoPath)
		if err != nil {
			cmd.Println(err)
			return
		}

		if len(workflows.Workflows) == 0 {
			cmd.Println("The repository does not have workflows.")
			return
		}

		for _, workflow := range workflows.Workflows {
			cmd.Printf("%s\n", workflow.GetName())
		}

	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoWorkflowCmd)

	repoListCmd.Flags().IntVarP(&size, "size", "s", 30, "Number of repositories to list")
	repoListCmd.Flags().BoolVar(&owned, "owned", false, "List only owned repos")
	repoListCmd.Flags().BoolVar(&followed, "followed", false, "List only followed repos")
	repoListCmd.MarkFlagsMutuallyExclusive("owned", "followed")
}

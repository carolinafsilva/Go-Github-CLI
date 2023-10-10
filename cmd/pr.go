package cmd

import (
	"github.com/carolinafsilva/go-github-cli/api"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	size   int
	status bool
)

var prCmd = &cobra.Command{
	Use:   "pr <command> [flags]",
	Short: "Get information about Github Pull Requests",
	Long: `The pr command in GG is designed to retrieve essential pull request information from GitHub.
	You can use this command to filter and display pull requests based on different criteria such as the author or repository`,
}

var prAuthorCmd = &cobra.Command{
	Use:   "author <username> [flags]",
	Short: "Get Pull Request information by author",
	Long:  `The author subcommand within the pr command allows you to filter and list pull requests on GitHub by the author's username. It provides a convenient way to find and inspect pull requests submitted by a specific user.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		author := args[0]

		prs, err := api.ListPRsByAuthor(author, size)
		if err != nil {
			cmd.Println(err)
			return
		}

		for i, pr := range prs {
			magenta := color.New(color.FgMagenta)
			cmd.Printf("%3d. ", i+1)
			magenta.Fprintf(cmd.OutOrStdout(), "%s ", *pr.CreatedAt)
			cmd.Printf("%s\n", *pr.Title)
		}
	},
}

var prRepoCmd = &cobra.Command{
	Use:   "repo <owner/repo> [flags]",
	Short: "Get Pull Request information by repository",
	Long:  `The repo subcommand within the pr command allows you to filter and list pull requests on GitHub by the repository name. You can use this subcommand to view pull requests associated with a particular repository.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]

		if status {
			prs, err := api.ListPRsByRepoWithStatus(repoPath, size)
			if err != nil {
				cmd.Println(err)
				return
			}

			for i, pr := range prs {
				statusColor := color.New(color.Bold)
				if pr.Status == "success" {
					statusColor.Add(color.FgGreen)
				} else if pr.Status == "pending" {
					statusColor.Add(color.FgYellow)
				} else {
					statusColor.Add(color.FgRed)
				}
				magenta := color.New(color.FgMagenta)

				cmd.Printf("%3d. ", i+1)
				magenta.Fprintf(cmd.OutOrStdout(), "%s", *pr.PR.CreatedAt)
				statusColor.Fprintf(cmd.OutOrStdout(), "  %s  ", pr.Status)
				cmd.Printf("%s\n", *pr.PR.Title)
			}
		} else {
			prs, err := api.ListPRsByRepo(repoPath, size)
			if err != nil {
				cmd.Println(err)
				return
			}

			for i, pr := range prs {
				magenta := color.New(color.FgMagenta)
				cmd.Printf("%3d. ", i+1)
				magenta.Fprintf(cmd.OutOrStdout(), "%s ", *pr.CreatedAt)
				cmd.Printf("%s\n", *pr.Title)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(prAuthorCmd)
	prCmd.AddCommand(prRepoCmd)

	prAuthorCmd.Flags().IntVarP(&size, "size", "S", 30, "Number of results to return")
	prRepoCmd.Flags().BoolVar(&status, "status", false, "Show the state of the PRs in the Workflow")
	prRepoCmd.Flags().IntVarP(&size, "size", "S", 30, "Number of results to return")
}

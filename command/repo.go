package command

import (
	"github.com/carthewd/c3/internal/awsclient"
	"github.com/carthewd/c3/pkg/codecommit"
	"github.com/carthewd/c3/util"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoListCmd)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
	Long: `Actions to view and manage CodeCommit repositories.
	
Returns a view of all repositories.`,
}

var repoListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "li"},
	Short:   "List all active CodeCommit repositories",
	RunE:    repoList,
}

func repoList(cmd *cobra.Command, args []string) error {
	profile, err := cmd.Flags().GetString("profile")
	c := awsclient.NewClient(profile)

	repoNames := codecommit.ListRepoNames(c)
	repoDetails := codecommit.GetRepoDetails(c, repoNames)

	util.PrintTable(repoDetails)
	return err
}

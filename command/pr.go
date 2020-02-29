package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/carthewd/c3/internal/awsclient"
	"github.com/carthewd/c3/pkg/codecommit"
	"github.com/carthewd/c3/pkg/gitconfig"
	"github.com/carthewd/c3/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(prCmd)
	prCmd.AddCommand(prListCmd)
	prCmd.AddCommand(prCOCmd)

	prListCmd.Flags().StringP("all", "a", "", "Show all <state> PRs for repository")
	prListCmd.Flags().StringP("state", "s", "open", "Show all <state> PRs for repository")
}

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "View PRs",
	Long: `Actions to view and manage CodeCommit pull requests.
	
A pull request can be supplied using the pull request ID, e.g., "321"`,
}

var prListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "li"},
	Short:   "List all PRs for a CodeCommit repository",
	RunE:    prList,
}

var prCOCmd = &cobra.Command{
	Use:     "checkout [pull request ID]",
	Aliases: []string{"co"},
	Short:   "Checkout a CodeCommit PR",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit PR number. ")
		}
		return nil
	},
	RunE: prCheckOut,
}

func prList(cmd *cobra.Command, args []string) error {
	var repo string
	err := error(nil)

	if repo, err = cmd.Flags().GetString("repo"); err == nil && repo == "" {
		repo, err = gitconfig.GetOrigin()
	}

	if err != nil || repo == "" {
		log.Fatal("No CodeCommit repository in current working directory.")
	}

	author, err := cmd.Flags().GetString("all")
	if err != nil {
		return err
	}

	state, err := cmd.Flags().GetString("state")
	if err != nil {
		return err
	}
	state = strings.ToUpper(state)

	c := awsclient.NewClient()

	prs := codecommit.ListPRs(c, repo, author, state)
	util.PrintTable(prs)

	return err
}

func prCheckOut(cmd *cobra.Command, args []string) error {
	err := error(nil)
	if len(args) == 0 {
		return err
	}

	c := awsclient.NewClient()
	pr, err := codecommit.GetPRDetails(c, args[0], "")

	gitcmd := []string{fmt.Sprintf(`{"fetch", "remote", %q}`, pr.SourceBranch)}
	gitcmd = append(gitcmd, pr.SourceBranch)

	o, _ := gitconfig.GitCmd(gitcmd...)
	fmt.Println(o)

	gitcmd = []string{"checkout"}
	gitcmd = append(gitcmd, pr.SourceBranch)

	o, _ = gitconfig.GitCmd(gitcmd...)

	fmt.Println(o)
	return err
}

func createPRURL(u string) string {
	//prURL := fmt.Sprintf("%s/codesuite/codecommit/repositories/%s/pull-requests/%s?region=%s", baseURL, repoName, prID, region)
	prURL := ""
	return prURL
}

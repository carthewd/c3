package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/carthewd/c3/internal/awsclient"
	"github.com/carthewd/c3/internal/data"
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
	prCmd.AddCommand(prDiffCmd)
	prCmd.AddCommand(prCreateCmd)

	prListCmd.Flags().StringP("author", "a", "", "Show <state> pull requests for repository by author (defaults to all)")
	prListCmd.Flags().StringP("state", "s", "open", "Show all <state> PRs for repository")
}

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create, view and checkout pull requests.",
	Long: `Actions to view and manage CodeCommit pull requests.
	
A pull request can be supplied using the pull request ID, e.g., "321"`,
}

var prListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "li"},
	Short:   "List all pull requests for a CodeCommit repository",
	RunE:    prList,
}

var prCOCmd = &cobra.Command{
	Use:     "checkout [pull request ID]",
	Aliases: []string{"co"},
	Short:   "Checkout a CodeCommit PR",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit pull request number. ")
		}
		return nil
	},
	RunE: prCheckOut,
}

var prDiffCmd = &cobra.Command{
	Use:     "diff [pull request ID]",
	Aliases: []string{"di"},
	Short:   "Show a diff for a given pull request",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit pull request number. ")
		}
		return nil
	},
	RunE: prDiff,
}

var prCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr"},
	Short:   "Create a pull request",
	RunE:    prCreate,
}

func prList(cmd *cobra.Command, args []string) error {
	repo, err := gitconfig.GetOrigin()

	if err != nil || repo == "" {
		log.Fatal("No CodeCommit repository in current working directory.")
	}

	author, err := cmd.Flags().GetString("author")
	if err != nil {
		return err
	}

	state, err := cmd.Flags().GetString("state")
	if err != nil {
		return err
	}
	state = strings.ToUpper(state)

	profile, err := cmd.Flags().GetString("profile")
	c := awsclient.NewClient(profile)

	prs := codecommit.ListPRs(c, repo, author, state)
	util.PrintTable(prs)

	return err
}

func prCheckOut(cmd *cobra.Command, args []string) error {
	err := error(nil)
	if len(args) == 0 {
		return err
	}

	profile, err := cmd.Flags().GetString("profile")
	c := awsclient.NewClient(profile)

	pr, err := codecommit.GetPRDetails(c, args[0], "")

	gitcmd := []string{fmt.Sprintf(`{"fetch", "remote", %q}`, pr.SourceBranch)}
	gitcmd = append(gitcmd, pr.SourceBranch)

	o, _ := gitconfig.GitCmd(gitcmd...)

	gitcmd = []string{"checkout"}
	gitcmd = append(gitcmd, pr.SourceBranch)

	o, _ = gitconfig.GitCmd(gitcmd...)

	fmt.Println(o)
	return err
}

func prDiff(cmd *cobra.Command, args []string) error {
	profile, err := cmd.Flags().GetString("profile")

	c := awsclient.NewClient(profile)

	pr, err := codecommit.GetPRCommits(c, args[0])

	_ = []string{fmt.Sprintf(`{"fetch", "remote", %q}`, pr.SourceBranch)}
	o, _ := gitconfig.GitCmd("diff", pr.DestCommit, pr.MergeCommit, "--color=always")

	fmt.Println(o)
	return err
}

func prCreate(cmd *cobra.Command, args []string) error {
	profile, err := cmd.Flags().GetString("profile")

	c := awsclient.NewClient(profile)

	newPR := data.NewPullRequest{}

	repo, err := gitconfig.GetOrigin()
	if err != nil || repo == "" {
		log.Fatal("No CodeCommit repository in current working directory.")
	}
	newPR.Repository = repo

	// Get current working branch
	srcBranch, _ := gitconfig.GitCmd("rev-parse", "--abbrev-ref", "HEAD")

	// Check branch exists in remote origin (i.e., change has been pushed)
	srcBranch = strings.Replace(srcBranch, "\n", "", 1)
	o, _ := gitconfig.GitCmd("ls-remote", "-q", "--heads", "origin", srcBranch)

	if o == "" {
		log.Fatal("No remote branch found - has your changed been pushed upstream?")
	}

	newPR.SourceRef = srcBranch

	o, _ = gitconfig.GitCmd("log", "-1", "--pretty=%B")
	newPR.Title = strings.Replace(o, "\n\n", "\n", 1)

	prTemplate := fmt.Sprintf(`%s
# ------------------------ >8 ------------------------
# Do not modify or remove the line above.
# Everything below it will be ignored.
`, newPR.Title)

	text, err := util.OpenInEditor(prTemplate)
	if err != nil {
		fmt.Println(err)
	}

	var str strings.Builder

	breader := bytes.NewReader(text)
	bufReader := bufio.NewReader(breader)

	tstr, _, _ := bufReader.ReadLine()
	newPR.Title = string(tstr)

	for {
		line, _, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		}

		if strings.Contains(string(line), "# ------------------------ >8 ------------------------") {
			break
		}
		str.WriteString(string(line) + "\n")
	}

	newPR.Description = str.String()

	result, err := codecommit.CreatePR(c, newPR)

	fmt.Println(util.CreatePullRequestURL(newPR.Repository, result))

	return err
}

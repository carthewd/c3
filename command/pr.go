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
	prCmd.AddCommand(prApprovalCmd)
	prCmd.AddCommand(prRevokeCmd)
	prCmd.AddCommand(prMergeCmd)

	prListCmd.Flags().StringP("author", "a", "", "Show <state> pull requests for repository by author (defaults to all)")
	prListCmd.Flags().StringP("state", "s", "open", "Show all <state> PRs for repository")

	prMergeCmd.Flags().BoolP("delete-branch", "d", true, "Delete the remote branch after a successful merge")
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

var prApprovalCmd = &cobra.Command{
	Use:     "approve [pull request ID]",
	Aliases: []string{"app", "appr", "a"},
	Short:   "Approve a pull request",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit pull request number. ")
		}
		return nil
	},
	RunE: prApprove,
}

var prRevokeCmd = &cobra.Command{
	Use:     "revoke [pull request ID]",
	Aliases: []string{"rev", "re", "r"},
	Short:   "Revoke a pull request",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit pull request number. ")
		}
		return nil
	},
	RunE: prRevoke,
}

var prCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr"},
	Short:   "Create a pull request",
	RunE:    prCreate,
}

var prMergeCmd = &cobra.Command{
	Use:	 "merge [pull request ID]",
	Aliases: []string{"m", "mr"},
	Short:	 "Merge a pull request",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit pull request number. ")
		}
		return nil
	},
	RunE:	 prMerge,
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

	_, _ = gitconfig.GitCmd("fetch", "origin", pr.SourceBranch)
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
		log.Fatal("No remote branch found - has your change been pushed upstream?")
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
	if strings.Contains(newPR.Description, "# Do not modify or remove the line above.") {
		log.Fatal("Empty description - not creating pull request.")
	} else if newPR.Description == "" {
		log.Fatal("Empty description - not creating pull request.")
	}

	result, err := codecommit.CreatePR(c, newPR)

	fmt.Println(util.CreatePullRequestURL(newPR.Repository, result))

	return err
}

func prApprove(cmd *cobra.Command, args []string) error {
	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return err
	}

	c := awsclient.NewClient(profile)

	err = codecommit.ApprovePR(c, args[0])

	if err != nil {
		return err
	}

	fmt.Printf("Pull request %s approved.\n", args[0])
	return err
}

func prRevoke(cmd *cobra.Command, args []string) error {
	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return err
	}

	c := awsclient.NewClient(profile)

	err = codecommit.RevokePR(c, args[0])

	if err != nil {
		return err
	}

	fmt.Printf("Pull request %s approval revoked.\n", args[0])
	return err
}

func prMerge(cmd *cobra.Command, args []string) error {
	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return err
	}

	c := awsclient.NewClient(profile)

	pr, err := codecommit.GetPRDetails(c, args[0], "")
	if err != nil {
		return err
	}

	repo, err := gitconfig.GetOrigin()
	if err != nil {
		return err
	}

	opts, _ := codecommit.MergeOptions(c, pr, repo)

	d, err := cmd.Flags().GetBool("delete-branch")
	if err != nil {
		return err
	}

	mergeInput := data.MergeInput{
		PRID: pr.ID,
		Repository: repo,
		SourceBranch: pr.SourceBranch,
		DeleteBranch: d,
	}

	if opts.FF {
		mergeInput.Type = "FF"
	} else if opts.ThreeWay {
		mergeInput.Type = "ThreeWay"
		log.Fatal("Only FastForward merge supported.")
	} else if opts.Squash {
		mergeInput.Type = "Squash"
		log.Fatal("Only FastForward merge supported.")
	} else {
		log.Fatal("No viable merge strategies - resolve conflicts and try again.")
	}

	_ = codecommit.Merge(c, mergeInput)

	fmt.Printf("Pull request %s merged successfully.\n", pr.ID)

	return nil
}

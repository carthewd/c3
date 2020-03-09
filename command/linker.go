package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/carthewd/c3/internal/data"
	"github.com/carthewd/c3/pkg/gitconfig"
	"github.com/carthewd/c3/util"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	RootCmd.AddCommand(cLinkCmd)
}

var cLinkCmd = &cobra.Command{
	Use:     "link",
	Aliases: []string{"l", "lnk"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a CodeCommit pull request number or file path")
		}
		return nil
	},
	RunE:  createLink,
	Short: "Generate AWS console links for CodeCommit objects.",
	Long: `Generate AWS console links for CodeCommit objects, such as file paths and pull requests. 
	
A pull request can be supplied using the pull request ID, e.g., "#321" or simply include the filename or directory.`,
}

func createLink(cmd *cobra.Command, args []string) error {
	var url string
	parsedArgs := strings.Split(args[0], ":")
	var err error

	if parsedArgs[0] == "pr" {
		repo, err := gitconfig.GetOrigin()

		if err != nil || repo == "" {
			log.Fatal("No CodeCommit repository in current working directory.")
		}
		url = util.CreatePullRequestURL(repo, parsedArgs[1])
	} else {
		info, err := os.Stat(parsedArgs[0])
		if os.IsNotExist(err) {
			log.Fatal("No such file or directory")
		}

		fp := data.Path{}
		if info.IsDir() {
			fp.PathType = "dir"
		} else {
			fp.PathType = "file"
		}

		fp.Path = args[0]
		url = util.CreatePathURL(fp)
	}

	fmt.Println(url)

	return err
}

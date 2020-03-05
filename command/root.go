package command

import (
	"github.com/spf13/cobra"
)

func init() {
}

var RootCmd = &cobra.Command{
	Use:   "c3",
	Short: "CodeCommit CLI",
	Long:  `Manage AWS CodeCommit workflows from the command line.`,
}

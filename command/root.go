package command

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().StringP("profile", "p", "", "AWS profile to use (from shared credentials file)")
}

var RootCmd = &cobra.Command{
	Use:   "c3",
	Short: "CodeCommit CLI",
	Long:  `Manage AWS CodeCommit workflows from the command line.`,
}

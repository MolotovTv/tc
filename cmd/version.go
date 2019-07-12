package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitCommit string

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gitCommit)
	},
}

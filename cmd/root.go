package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "tc",
	Short: "TeamCity",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

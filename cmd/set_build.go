package cmd

import (
	"log"

	"github.com/molotovtv/tc/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(setBuildCmd)
}

var setBuildCmd = &cobra.Command{
	Use:   "set-build",
	Short: "Set a build id",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SetBuildTypeID(args[0], args[1]); err != nil {
			log.Fatal(err)
		}
	},
}

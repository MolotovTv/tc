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
		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		err = c.SetBuildID(projectName(), args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
	},
}

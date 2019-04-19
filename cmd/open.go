package cmd

import (
	"log"

	"github.com/molotovtv/tc/internal/config"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "op",
	Short: "Open the teamcity website",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}
		buildID, err := config.BuildID(env)
		if err != nil {
			log.Fatal(err)
		}
		open.Run(c.URL + "/viewType.html?buildTypeId=" + buildID)
	},
}

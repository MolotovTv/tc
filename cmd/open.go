package cmd

import (
	"log"

	"github.com/molotovtv/tc/tc"
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
		c, err := tc.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}
		buildTypeID, err := c.BuildTypeID(env)
		if err != nil {
			log.Fatal(err)
		}
		open.Run(c.URL + "/viewType.html?buildTypeId=" + buildTypeID)
	},
}

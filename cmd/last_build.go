package cmd

import (
	"fmt"
	"log"

	"github.com/molotovtv/tc/tc"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(lastBuildCmd)
}

var lastBuildCmd = &cobra.Command{
	Use:   "last-build",
	Short: "Get last build info",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		fmt.Println("Env:", env)

		c, err := tc.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		buildTypeID, err := c.BuildTypeID(env)
		if err != nil {
			log.Fatal(err)
		}

		res, err := tc.LastBuild(c, buildTypeID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", res)
	},
}

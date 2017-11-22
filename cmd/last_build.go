package cmd

import (
	"fmt"
	"log"

	"github.com/aestek/tc/tc"

	"github.com/aestek/tc/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(lastBuildCmd)
}

var lastBuildCmd = &cobra.Command{
	Use:   "last-build",
	Short: "get last build info",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		fmt.Println("Env:", env)

		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		buildID := c.BuildIDPrompt(projectName(), env)

		res, err := tc.LastBuild(c, buildID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", res)
	},
}

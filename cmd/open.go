package cmd

import (
	"log"

	"github.com/aestek/tc/internal/config"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "op",
	Short: "Open a project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]

		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		buildID := c.BuildIDPrompt(projectName(), env)

		open.Run(c.URL + "/viewType.html?buildTypeId=" + buildID)
	},
}

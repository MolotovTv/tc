package cmd

import (
	"fmt"
	"log"

	"github.com/aestek/tc/internal/tc"

	"github.com/aestek/tc/internal/config"
	"github.com/aestek/tc/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a build",
	Run: func(cmd *cobra.Command, args []string) {
		branch, err := git.Branch()
		if err != nil {
			log.Fatal(err)
		}

		local, err := config.Local()
		if err != nil {
			log.Fatal(err)
		}

		global, err := config.Global()
		if err != nil {
			log.Fatal(err)
		}

		buildID := local.BuildIDPromp(config.EnvStag)

		res, err := tc.RunBranch(global, buildID, branch)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
	},
}

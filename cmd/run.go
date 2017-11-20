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
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		fmt.Println("Env:", env)

		branch, err := git.Branch()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Branch:", branch)

		local, err := config.Local()
		if err != nil {
			log.Fatal(err)
		}

		global, err := config.Global()
		if err != nil {
			log.Fatal(err)
		}

		buildID := local.BuildIDPromp(env)

		res, err := tc.RunBranch(global, buildID, branch)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
	},
}

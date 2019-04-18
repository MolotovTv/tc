package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/molotovtv/tc/tc"

	"github.com/molotovtv/tc/internal/config"
	"github.com/molotovtv/tc/internal/git"
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

		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		buildID := c.BuildIDPrompt(projectName(), env)

		lastBuild, err := tc.LastBuild(c, buildID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = tc.RunBranch(c, buildID, branch)
		if err != nil {
			log.Fatal(err)
		}

		for {
			time.Sleep(time.Second)

			build, err := tc.LastBuild(c, buildID)
			if err != nil {
				log.Fatal(err)
			}

			if build.ID != lastBuild.ID {
				break
			}
		}

		buildStatus(c, buildID)
	},
}

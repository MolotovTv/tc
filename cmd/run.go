package cmd

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/molotovtv/tc/tc"

	"github.com/molotovtv/tc/internal/config"
	"github.com/molotovtv/tc/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(runCmd)
}

func renameBranchForProd(branchOrigin string) string {
	re, err := regexp.Compile("^release\\-([0-9]+(\\.[0-9]+)*)$")
	if err != nil {
		return ""
	}
	matches := re.FindAllStringSubmatch(branchOrigin, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return ""
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
		if env == "prod" {
			prodBranch := renameBranchForProd(branch)
			if branch == "" {
				log.Fatal(fmt.Errorf("could not remap branch %s to a correct prod branch name", branch))
			}
			branch = prodBranch
		}

		fmt.Println("Branch:", branch)

		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		buildID, err := config.BuildID(env)
		if err != nil {
			log.Fatal(err)
		}

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

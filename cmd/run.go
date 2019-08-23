package cmd

import (
	"fmt"
	"log"
	"regexp"

	"github.com/fatih/color"
	"github.com/molotovtv/tc/tc"

	"github.com/molotovtv/tc/internal/config"
	"github.com/molotovtv/tc/internal/git"
	"github.com/spf13/cobra"
)

var (
	runSilent bool
	regex     = regexp.MustCompile("^release\\-([0-9]+(\\.[0-9]+)*)$")
)

func init() {
	runCmd.PersistentFlags().BoolVarP(&runSilent, "silent", "s", false, "Silent")
	RootCmd.AddCommand(runCmd)
}

func renameBranchForProd(branchOrigin string) string {
	matches := regex.FindAllStringSubmatch(branchOrigin, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return ""
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Deploy a version",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]

		fmt.Println("Env:", env)
		var tag string
		branch, err := git.CurrentBranch()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if env == "prod" || env == "test" {
			if len(args) > 1 {
				tag = args[1]
			} else {
				tags, err := git.CurrentTags()
				if err != nil {
					log.Fatalf("%+v", err)
					return
				}
				if len(tags) == 0 {
					log.Fatalf("No tag found in HEAD. precise which tag to run.")
					return
				} else if len(tags) > 1 {
					log.Fatalf("more than one tag available, precise which one to package: %+v", tags)
					return
				} else {
					tag = tags[0]
				}
			}

			if !tagRegex.Match([]byte(tag)) {
				log.Fatalf("tag %s does not validate, tags for prod must be vx.y.z", tag)
				return
			}
			branch = tag
		}

		c, err := config.Load()
		if err != nil {
			log.Fatalf("%+v", err)
		}

		buildTypeID, err := config.BuildTypeID(env)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if env == "prod" {
			refInfos, err := git.ShowRef(tag)
			if err != nil {
				log.Fatalf("%+v", err)
			}
			branch = branch[1:len(branch)]
			fmt.Printf("Will deploy version (%s) with this commit to %s\n------------------------------\n%s\n------------------------------\nContinue?", color.New(color.FgGreen).SprintFunc()(branch), color.New(color.Bold).SprintFunc()("prod"), refInfos)
			var ok string
			fmt.Scanln(&ok)
		}
		fmt.Printf("Deploying version %s on %s...\n", color.New(color.FgGreen).SprintFunc()(branch), color.New(color.FgGreen).SprintFunc()(env))

		buildID, err := tc.RunBranch(c, buildTypeID, branch)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		if runSilent {
			return
		}

		buildStatus(c, buildID)
	},
}

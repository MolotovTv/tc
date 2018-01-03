package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/0xAX/notificator"
	"github.com/aestek/tc/internal/config"
	"github.com/aestek/tc/tc"
	"github.com/cheggaaa/pb"
	"github.com/spf13/cobra"
)

func buildStatus(c *config.Config, buildID string) {
	notify := notificator.New(notificator.Options{
		AppName: "TeamCity",
	})

	build, err := tc.LastBuild(c, buildID)
	if err != nil {
		log.Fatal(err)
	}

	if build.State == tc.BuildStatusFinished {
		return
	}

	bar := pb.StartNew(100)
	defer func() {
		bar.FinishPrint("done.")
		notify.Push("Build finished", build.BranchName, "", notificator.UR_NORMAL)
	}()

	for {
		build, err := tc.LastBuild(c, buildID)
		if err != nil {
			log.Fatal(err)
		}

		if build.State == tc.BuildStatusFinished {
			return
		}

		bar.Set(int(build.PercentageComplete))
		time.Sleep(time.Second)
	}
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "st",
	Short: "Build status",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		fmt.Println("Env:", env)

		c, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		buildID := c.BuildIDPrompt(projectName(), env)
		buildStatus(c, buildID)
	},
}

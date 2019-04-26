package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/molotovtv/tc/tc"

	"github.com/molotovtv/tc/internal/config"
	"github.com/spf13/cobra"
)

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

		buildID, err := config.BuildID(env)
		if err != nil {
			log.Fatal(err)
		}
		buildStatus(c, buildID)
	},
}

func buildStatus(c config.Config, buildID string) {
	bar := pb.StartNew(100)
	var build *tc.Build
	for {
		var err error
		build, err = tc.LastBuild(c, buildID)
		if err != nil {
			log.Fatal(err)
		}

		if build.State == tc.BuildStatusFinished {
			break
		}

		bar.Set(int(build.PercentageComplete))
		time.Sleep(time.Second)
	}
	bar.FinishPrint("done.")
	fmt.Printf("last build:\n%+v", build)
}

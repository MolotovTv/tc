package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/molotovtv/tc/internal/config"
	"github.com/molotovtv/tc/tc"
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

		buildTypeID, err := config.BuildTypeID(env)
		if err != nil {
			log.Fatal(err)
		}

		build, err := tc.LastBuild(c, buildTypeID)
		if err != nil {
			log.Fatal(err)
		}
		buildStatus(c, build.ID)
	},
}

func buildStatus(c config.Config, buildID int) {
	bar := pb.StartNew(100)
	var build tc.DetailedBuild
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	buildInterrupted := false
L:
	for {
		var err error
		build, err = tc.GetBuild(c, buildID)
		if err != nil {
			log.Fatal(err)
		}

		if build.State == tc.BuildStatusFinished {
			break
		}

		bar.Set(int(build.PercentageComplete))
		select {
		case <-signalChan:
			buildInterrupted = true
			break L
		case <-time.After(time.Second * 3):
		}
	}
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	bar.Finish()
	if buildInterrupted {
		if err := tc.CancelBuild(c, build.ID); err != nil {
			log.Fatal(err)
		}
		color.Magenta("Build cancelled!")
		return
	}
	if build.Status != tc.BuildStatusSuccess {
		color.Red("Build failed!")
		return
	}
	color.Green("Build succeeded!")
}

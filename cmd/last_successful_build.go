package cmd

import (
	"fmt"
	"log"

	"github.com/molotovtv/tc/tc"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(buildsCmd)
}

var buildsCmd = &cobra.Command{
	Use:   "last-build-success",
	Short: "Get successful builds by branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("builds-success [env] [branch]\n")
			return
		}
		env := args[0]
		branch := args[1]
		fmt.Printf("Env:%s ; Branch:%s\n", env, branch)

		c, err := tc.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		buildTypeID, err := c.BuildTypeID(env)
		if err != nil {
			log.Fatal(err)
		}

		res, err := tc.LastBuildSuccessByBranch(c, buildTypeID, branch)
		if err != nil { 
			log.Fatal(err) 
		}

		fmt.Printf("%+v\n", res)
	},
}

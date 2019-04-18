package cmd

import (
	"log"

	"github.com/molotovtv/tc/server"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run ther server",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Run(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

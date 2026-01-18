package main

import (
	"fmt"
	"log"

	"efrc/internal/client"

	"github.com/spf13/cobra"
)

var peersCmd = &cobra.Command{
	Use:   "peers",
	Short: "List all devices in the network",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := client.LoadConfig()
		if err != nil {
			log.Fatalf("No config found. Run 'efrc init' first: %v", err)
		}

		if config.ServerURL == "" {
			log.Fatal("Not logged in. Run 'efrc login' first")
		}

		// TODO: Fetch peers from server
		fmt.Println("Peers: (not implemented yet)")
	},
}

func init() {
	rootCmd.AddCommand(peersCmd)
}

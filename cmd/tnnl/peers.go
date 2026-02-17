package main

import (
	"fmt"
	"log"

	"github.com/eggfriedrice24/tnnl/internal/client"

	"github.com/spf13/cobra"
)

var peersCmd = &cobra.Command{
	Use:   "peers",
	Short: "List all devices in the network",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := client.LoadConfig()
		if err != nil {
			log.Fatalf("No config found. Run 'tnnl init' first: %v", err)
		}

		if config.ServerURL == "" {
			log.Fatal("Not logged in. Run 'tnnl login' first")
		}

		api := client.NewAPIClient(config.ServerURL, config.NetworkKey)
		peers, err := api.GetPeers()
		if err != nil {
			log.Fatalf("Failed to get peers: %v", err)
		}

		if len(peers) == 0 {
			fmt.Println("No peers found")
			return
		}

		fmt.Printf("%-20s %-15s %-10s\n", "NAME", "VIRTUAL IP", "ONLINE")
		fmt.Println("--------------------------------------------")
		for _, p := range peers {
			fmt.Printf("%-20s %-15s %-10v\n", p.Name, p.VirtualIP, p.Online)
		}
	},
}

func init() {
	rootCmd.AddCommand(peersCmd)
}

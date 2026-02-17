package main

import (
	"fmt"
	"log"

	"github.com/eggfriedrice24/tnnl/internal/client"

	"github.com/spf13/cobra"
)

var serverURL string

var loginCmd = &cobra.Command{
	Use:   "login [network-key]",
	Short: "Configure server and network key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := client.LoadConfig()
		if err != nil {
			log.Fatalf("No config found. Run 'tnnl init' first: %v", err)
		}

		config.NetworkKey = args[0]
		config.ServerURL = serverURL

		if err := config.Save(); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}

		fmt.Println("Login configured!")
		fmt.Printf("Server: %s \n", config.ServerURL)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&serverURL, "server", "s", "", "Server URL")
	loginCmd.MarkFlagRequired("server")
}

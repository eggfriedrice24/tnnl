package main

import (
	"fmt"
	"log"

	"github.com/eggfriedrice24/tnnl/internal/client"
	"github.com/eggfriedrice24/tnnl/internal/protocol"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register this device with the server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := client.LoadConfig()
		if err != nil {
			log.Fatalf("No config found. Run 'tnnl init' first: %v", err)
		}

		if config.ServerURL == "" {
			log.Fatal("Not logged in. Run 'tnnl login' first")
		}

		api := client.NewAPIClient(config.ServerURL, config.NetworkKey)

		req := &protocol.RegisterRequest{
			Name:      config.DeviceName,
			PublicKey: config.PublicKey,
			Endpoints: []string{}, // TODO: discover local endpoints
		}

		resp, err := api.Register(req)
		if err != nil {
			log.Fatalf("Registration failed: %v", err)
		}

		// Save the response to config
		config.DeviceID = resp.DeviceID
		config.VirtualIP = resp.VirtualIP
		if err := config.Save(); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}

		fmt.Println("Registered!")
		fmt.Printf("Device ID:  %s\n", resp.DeviceID)
		fmt.Printf("Virtual IP: %s\n", resp.VirtualIP)
		fmt.Printf("Peers:      %d\n", len(resp.Peers))
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

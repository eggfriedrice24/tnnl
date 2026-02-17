package main

import (
	"fmt"
	"log"

	"github.com/eggfriedrice24/tnnl/internal/client"

	"github.com/spf13/cobra"
)

var deviceName string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize device and generate WireGuard keys",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := client.InitConfig(deviceName)
		if err != nil {
			log.Fatalf("Failed to initialize: %v", err)
		}

		fmt.Println("Device initialized!")
		fmt.Printf("Name:       %s\n", config.DeviceName)
		fmt.Printf("Public Key: %s\n", config.PublicKey)
		fmt.Printf("Config:     ~/.tnnl/config.json\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&deviceName, "name", "n", "", "Device name")
	initCmd.MarkFlagRequired("name")
}

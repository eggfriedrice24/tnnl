package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eggfriedrice24/tnnl/internal/server"

	"github.com/spf13/cobra"
)

var (
	port       int
	networkKey string
)

var rootCmd = &cobra.Command{
	Use:   "tnnl-server",
	Short: "Control server for tnnl mesh VPN",
	Run: func(cmd *cobra.Command, args []string) {
		if networkKey == "" {
			log.Fatal("--network-key is required")
		}

		store := server.NewStore()
		api := server.NewAPI(store, networkKey)
		handler := api.SetupRoutes()

		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Server starting on %s\n", addr)
		log.Fatal(http.ListenAndServe(addr, handler))
	},
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 2424, "Server port")
	rootCmd.Flags().StringVarP(&networkKey, "network-key", "k", "", "Networkauthentication key")
	rootCmd.MarkFlagRequired("network-key")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

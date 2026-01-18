package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"efrc/internal/server"
)

func main() {
	port := flag.Int("port", 2424, "Server port")
	networkKey := flag.String("network-key", "", "Network authentication key")
	flag.Parse()

	if *networkKey == "" {
		log.Fatal("--network-key is required")
	}

	store := server.NewStore()
	api := server.NewAPI(store, *networkKey)
	handler := api.SetupRoutes()

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

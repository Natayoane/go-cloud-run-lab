package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ntayoane/go-cloud-run-lab/configs"
	"github.com/ntayoane/go-cloud-run-lab/internal"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/temperature", internal.HandleTemperatureRequest(*config))
	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

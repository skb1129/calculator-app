package main

import (
	"log"
	"os"

	"github.com/skb1129/calculator-app/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		// Default to relative frontend/dist if present
		if _, err := os.Stat("frontend/dist"); err == nil {
			staticDir = "frontend/dist"
		}
	}

	router := api.SetupRouter(staticDir)

	log.Printf("Starting calculator server on port :%s (static files: %q)", port, staticDir)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}

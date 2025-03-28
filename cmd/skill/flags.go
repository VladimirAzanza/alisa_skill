package main

import (
	"flag"
	"os"
)

var flagRunAddr string

// parseFlags initializes and parses command-line flags for the application.
// It sets the default server address to ":8080" and allows customization
// via the "-a" flag (e.g., "-a=:9090" to change the listen address).
// The parsed value is stored in the global variable `flagRunAddr`.
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
}

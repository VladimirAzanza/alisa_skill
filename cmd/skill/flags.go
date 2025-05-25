package main

import (
	"flag"
	"os"
)

var (
	flagRunAddr     string
	flagLogLevel    string
	flagDatabaseURI string
)

// parseFlags initializes and parses command-line flags for the application.
// It sets the default server address to ":8080" and allows customization
// via the "-a" flag (e.g., "-a=:9090" to change the listen address).
// The parsed value is stored in the global variable `flagRunAddr`.
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "debug", "log level")
	flag.StringVar(&flagDatabaseURI, "d", "", "database URI")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		flagDatabaseURI = envDatabaseURI
	}
}

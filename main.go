package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uguremirmustafa/inventory/api"
	"github.com/uguremirmustafa/inventory/logging"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a file logger instance
	logger := logging.NewFileLogger("app.log")

	// Use the logger to log messages
	logger.Infof("Starting application...")

	ctx := context.Background()
	if err := api.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uguremirmustafa/inventory/api"
	logging "github.com/uguremirmustafa/inventory/log"
)

func main() {
	l := logging.NewFileLogger("app.log")
	l.Infof("Starting application...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	if err := api.Run(ctx, l); err != nil {
		l.Errorf("Error running application: %v", err)
		os.Exit(1)
	}
}

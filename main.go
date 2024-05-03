package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uguremirmustafa/inventory/api"
)

func init() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logHandler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "date"
				a.Value = slog.AnyValue(time.Now().Format("2006/01/02 15:04:05"))
			}
			return a
		},
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func main() {
	slog.Info("Starting application...")

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	ctx := context.Background()
	if err := api.Run(ctx); err != nil {
		slog.Error("Error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

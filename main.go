package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/uguremirmustafa/inventory/api"
	"github.com/uguremirmustafa/inventory/internal/config"
)

func init() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logHandler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
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

	err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config.json file")
	}

	ctx := context.Background()
	if err := api.Run(ctx); err != nil {
		slog.Error("Error running application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

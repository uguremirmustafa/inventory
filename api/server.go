package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/uguremirmustafa/inventory/db"
)

func Run(ctx context.Context) error {
	pgDB, err := NewPostgresDB("postgres://anomy:secret@localhost:5432/inventory?sslmode=disable")
	if err != nil {
		panic(err)
	}
	q := db.New(pgDB)

	// w.Write([]byte(args[1]))
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	c := NewConfig()
	srv := NewServer(c, q)
	httpServer := &http.Server{
		Addr:    ":9000",
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown (thanks Alessandro Rosetti)
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewServer(c *Config, q *db.Queries) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, q, c)
	return mux
}

func addRoutes(mux *http.ServeMux, q *db.Queries, c *Config) {
	mux.Handle("POST /v1/users", handleGreet(q, c))
	mux.Handle("GET /v1/auth/login", handleLoginGoogle(c))
	mux.Handle("GET /v1/auth/callback", handleCallbackGoogle(q, c))
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (*T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return &v, fmt.Errorf("decode json: %w", err)
	}
	return &v, nil
}

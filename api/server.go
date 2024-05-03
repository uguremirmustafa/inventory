package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/justinas/alice"
	"github.com/uguremirmustafa/inventory/db"
)

func Run(ctx context.Context) error {

	pgDB, err := NewPostgresDB("postgres://anomy:secret@localhost:5432/inventory?sslmode=disable")
	if err != nil {
		panic(err)
	}
	q := db.New(pgDB)
	defer pgDB.Close()

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
		slog.Info("server is listening", slog.String("port", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error listening and serving", slog.String("error", err.Error()))
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("error shutting down http server", slog.String("error", err.Error()))
		}
		slog.Info("Server stopped gracefully")
	}()
	wg.Wait()
	return nil
}

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
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

type Middleware = func(http.Handler) http.Handler

// LoggingMiddleware logs information about each incoming request.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("Incoming request", slog.Group("request", slog.String("method", r.Method), slog.String("path", r.URL.Path)))
		next.ServeHTTP(w, r)
		slog.Info("Request time", slog.Duration("time", time.Since(start)))
	})
}

func addRoutes(mux *http.ServeMux, q *db.Queries, c *Config) {
	chain := alice.New(logMiddleware)
	authChain := alice.New(logMiddleware, authMiddleware(c))

	mux.Handle("GET /", chain.Then(handleHome()))
	mux.Handle("POST /v1/users", chain.Then(handleGreet(q, c)))
	mux.Handle("GET /v1/auth/login", chain.Then(handleLoginGoogle(c)))
	mux.Handle("GET /v1/auth/callback", chain.Then(handleCallbackGoogle(q, c)))
	mux.Handle("GET /v1/me", authChain.Then(handleMe(q, c)))
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

func handleHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		<div>
			<a href="/v1/auth/login">Login with Google</a>
			<a href="/v1/me">see me</a>
		</div>`)
	})
}

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
	"github.com/rs/cors"
	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/internal/config"
)

func Run(ctx context.Context) error {
	c := config.GetConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Database.Hostname,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
	)
	pgDB, err := NewPostgresDB(psqlInfo)
	if err != nil {
		panic(err)
	}
	q := db.New(pgDB)
	defer pgDB.Close()

	// w.Write([]byte(args[1]))
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv := NewServer(q, pgDB)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", c.PORT),
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

func NewServer(q *db.Queries, db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, q, db)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "X-Auth-Token"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie"},
	})

	return c.Handler(mux)
}

type Middleware = func(http.Handler) http.Handler

// LoggingMiddleware logs information about each incoming request.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Debug("Incoming request", slog.Group("request", slog.String("method", r.Method), slog.String("path", r.URL.Path)))
		next.ServeHTTP(w, r)
		slog.Debug("Request time", slog.Duration("took", time.Since(start)))
	})
}

func addRoutes(mux *http.ServeMux, q *db.Queries, db *sql.DB) {
	chain := alice.New(logMiddleware)
	authChain := alice.New(logMiddleware, authMiddleware())

	authService := NewAuthService(q, db)
	mux.Handle("POST /v1/auth/login", chain.Then(Make(authService.HandleLogin)))

	// deprecate
	// mux.Handle("GET /v1/auth/login", chain.Then(handleLoginGoogle()))
	// mux.Handle("GET /v1/auth/callback", chain.Then(handleCallbackGoogle(q)))
	mux.Handle("GET /v1/me", authChain.Then(handleMe(q)))

	itemTypeService := NewItemTypeService(q)
	mux.Handle("GET /v1/item-type", authChain.Then(Make(itemTypeService.HandleListItemType)))
	mux.Handle("POST /v1/item-type", authChain.Then(Make(itemTypeService.HandleCreateItemType)))

	mux.Handle("GET /v1/manufacturer", authChain.Then(handleListManufacturer(q)))
	mux.Handle("GET /v1/location", authChain.Then(handleListLocation(q)))

	itemService := NewItemService(q, db)
	mux.Handle("GET /v1/item", authChain.Then(Make(itemService.HandleListUserItem)))
	mux.Handle("POST /v1/item", authChain.Then(Make(itemService.HandleInsertUserItem)))
}

func encode(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func writeJson(w http.ResponseWriter, status int, v interface{}) error {
	return encode(w, status, Result{
		StatusCode: status,
		Data:       v,
	})
}

func decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return InvalidJSON()
	}
	return nil
}

type Result struct {
	StatusCode int `json:"statusCode"`
	Data       any `json:"data"`
}

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

func NotFound() APIError {
	return NewAPIError(http.StatusNotFound, fmt.Errorf("no items found"))
}

func NotAuthorized() APIError {
	return NewAPIError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid JSON in request body"))
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				// encode(w, apiErr.StatusCode, apiErr)
				writeJson(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				writeJson(w, http.StatusInternalServerError, errResp)
			}
			slog.Error("HTTP API Error", slog.String("err", err.Error()), slog.String("path", r.URL.Path))
		}
	}
}

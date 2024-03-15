// Package server implements one JSON api endpoint, api/v1/items, for the purpose
// of comparing its performance vs Rails and Sinatra
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charliemcelfresh/go-items-api/internal/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Item struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type Repository interface {
	GetItems(ctx context.Context, page int) ([]Item, error)
}

type server struct {
	Repository Repository
}

func NewServer() server {
	r := NewRepository()
	return server{
		Repository: r,
	}
}

// Run runs a Go stdlib http server.
// It uses a Chi router, see https://github.com/go-chi/chi,
// with Chi's logging middleware, and some custom middlewares.
func (s server) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(AddUserIDToContext)
	r.Use(AddContentTypeToResponse)
	r.Get("/api/v1/items", s.GetItems)

	httpServer := &http.Server{
		Addr:    ":3001",
		Handler: r,
	}
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server listening on :3001")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("ListenAndServe: %v", err))
		}
	}()

	<-stop

	config.GetLogger().Info("Shutting down httpServer ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Server graceful shutdown failed: %v", err))
	}

	config.GetLogger().Info("Server shutdown")
}

// GetItems retrieves items for the current user, identified by their
// X-User-Id header, with an offset determined by the page=? URL param
func (s server) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pageAsInt := getPage(r)
	items, err := s.Repository.GetItems(ctx, pageAsInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	jsonItems, err := json.Marshal(items)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(jsonItems)
}

// getPage retrieves the page param from the URL, if it
// exists. If it does not, page is set to 0
func getPage(r *http.Request) int {
	s := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return pageInt
}

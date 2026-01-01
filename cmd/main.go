package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	todoApi "github.com/sborsh1kmusora/todo/internal/api/todo"
	todoRepo "github.com/sborsh1kmusora/todo/internal/repository/todo"
	todoServ "github.com/sborsh1kmusora/todo/internal/service/todo"
)

const (
	serverAddr = "localhost:8080"

	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second
)

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	mux := http.NewServeMux()

	repo := todoRepo.New()
	service := todoServ.New(&repo)
	api := todoApi.New(&service, log)

	mux.HandleFunc("/todos", api.Todos)
	mux.HandleFunc("/todos/{id}", api.TodoById)

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Info("server is running on address", slog.String("address", server.Addr))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Info("failed to run server", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Info("failed to stop server", slog.String("error", err.Error()))
	}

	log.Info("server stopped gracefully")
}

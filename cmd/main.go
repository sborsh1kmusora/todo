package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	todoApi "github.com/sborsh1kmusora/todo/internal/api/todo"
	todoRepo "github.com/sborsh1kmusora/todo/internal/repository/todo"
	todoServ "github.com/sborsh1kmusora/todo/internal/service/todo"
)

const (
	configPath = ".env"

	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second
)

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	if err := godotenv.Load(configPath); err != nil {
		log.Error("failed to load .env file", slog.String("error", err.Error()))
		return
	}

	mux := http.NewServeMux()

	repo := todoRepo.New()
	service := todoServ.New(&repo)
	api := todoApi.New(&service)

	mux.HandleFunc("/todos", api.Create)

	server := &http.Server{
		Addr:              serverAddr(),
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

func serverAddr() string {
	return net.JoinHostPort(os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
}

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sjana7797/students/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup

	// router setup
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// server setup
	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}
	slog.Info("Server started", slog.String("address", cfg.HttpServer.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("Shutting down Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

}

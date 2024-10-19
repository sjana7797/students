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
	"github.com/sjana7797/students/internal/http/handlers/student"
	"github.com/sjana7797/students/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialed", slog.String("env", cfg.Env), slog.String("version", cfg.Version))

	// router setup
	router := http.NewServeMux()

	// students route
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students", student.GetStudents(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetStudent(storage))

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

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

}

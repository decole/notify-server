package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"notify-server/internal/config"
	"notify-server/internal/http-server/handlers/save"
	mvLogger "notify-server/internal/http-server/middleware"
	"notify-server/internal/lib/sl"
	"notify-server/internal/storage/postgres"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	_ = os.Setenv("CONFIG_PATH", "./config/config.yaml") // todo fix it

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("environment", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mvLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/notify", save.New(log, storage))

	log.Info("server started", slog.String("address", cfg.Server.Address))

	srv := &http.Server{
		Addr:              cfg.Server.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.Server.Timeout,
		WriteTimeout:      cfg.Server.Timeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/babenow/url_shortener/intrernal/config"
	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
	"github.com/babenow/url_shortener/intrernal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.Instance()

	logger := setupLogger(cfg.Env)
	logger.Info("starting application", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	storage, err := sqlite.New(ctx)
	if err != nil {
		logger.Error("can not initialize storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO: init router
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID, middleware.RealIP)
	router.Use(middleware.Logger)
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case config.EnvProd:

		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatalf("[ERROR]: wrong environment")
	}

	return logger
}

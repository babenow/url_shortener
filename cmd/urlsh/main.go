package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/babenow/url_shortener/intrernal/config"
)

func main() {
	cfg := config.Instance()

	logger := setupLogger(cfg.Env)
	logger.Info("starting application", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// TODO: init db

	// TODO: init router

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

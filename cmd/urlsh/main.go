package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/babenow/url_shortener/intrernal/config"
	"github.com/babenow/url_shortener/intrernal/http-server/handlers/redirect"
	"github.com/babenow/url_shortener/intrernal/http-server/handlers/url/save"
	"github.com/babenow/url_shortener/intrernal/http-server/middleware/logger"
	"github.com/babenow/url_shortener/intrernal/lib/logger/handlers/slogpretty"
	"github.com/babenow/url_shortener/intrernal/lib/logger/sl"
	"github.com/babenow/url_shortener/intrernal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// TODO: move timeout to config
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Instance()

	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(ctx, log)
	if err != nil {
		log.Error("can not initialize storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID, middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) // TODO: привязка к роутеру chi

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HttpServer.User: cfg.HttpServer.Password,
		}))
		r.Post("/", save.New(log, storage.UrlStorage()))
		// TODO: delete url
	})

	router.Get("/{alias}", redirect.New(ctx, log, storage.UrlStorage()))

	// starting server
	log.Info("Starting server", slog.String("address", cfg.HttpServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.ErrorContext(ctx, "failed to start sertver", sl.Err(err))
		}
	}()

	log.Info("server started")

	<-done
	log.Error("stopping server")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case config.EnvLocal:
		logger = setupPrettyLogger()
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

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

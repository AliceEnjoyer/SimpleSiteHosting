package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/AliceEnjoyer/SimpleSiteHosting/internal/config"
	"github.com/AliceEnjoyer/SimpleSiteHosting/internal/http-server/handlers/getPage"
	"github.com/AliceEnjoyer/SimpleSiteHosting/internal/http-server/middleware/logger"
	"github.com/AliceEnjoyer/SimpleSiteHosting/lib/logger/sl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ConfigPath := flag.String(
		"ConfigPath",
		"",
		"",
	)

	flag.Parse()
	*ConfigPath = "./config/config.yaml"
	if *ConfigPath == "" {
		log.Fatal("ConfigPath is not specified")
	}

	cnfg := config.MustLoad(*ConfigPath)
	log := setupLogger(cnfg.Env)
	log.Info("starting hosting", slog.String("env", cnfg.Env))
	log.Debug("debug messages are enabled")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/{pageName}", getPage.New(log, cnfg.PagesPath))

	srv := &http.Server{
		Addr:         cnfg.Address,
		Handler:      router,
		ReadTimeout:  cnfg.HTTPServer.Timeout,
		WriteTimeout: cnfg.HTTPServer.Timeout,
		IdleTimeout:  cnfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stoped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
		)
	}

	return log
}

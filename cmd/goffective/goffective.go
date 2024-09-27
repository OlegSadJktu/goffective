package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/OlegSadJktu/goffective/internal/config"
	mid "github.com/OlegSadJktu/goffective/internal/httpserver/middleware"
	"github.com/OlegSadJktu/goffective/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

var (
	env = flag.String("env", ".env", "Configuration file")
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}

func parseDbUrl(env *config.Config) string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		env.Postgres.User,
		env.Postgres.Password,
		env.Postgres.Host,
		env.Postgres.Port,
		env.Postgres.Database,
	)
}

func main() {
	flag.Parse()
	cfg := config.MustLoad(*env)
	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	url := parseDbUrl(cfg)
	_, err := storage.New(url)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(mid.NewHttpLogger(log))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Route("/songs", func(r chi.Router) {
	})

	err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.Server.Port), router)
	slog.Error(err.Error())

}

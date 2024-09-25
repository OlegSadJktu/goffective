package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/OlegSadJktu/goffective/internal/config"
	"github.com/OlegSadJktu/goffective/internal/storage"
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

}

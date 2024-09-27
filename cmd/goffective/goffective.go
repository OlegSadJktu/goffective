package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/OlegSadJktu/goffective/internal/config"
	mid "github.com/OlegSadJktu/goffective/internal/httpserver/middleware"
	"github.com/OlegSadJktu/goffective/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

var (
	env            = flag.String("env", ".env", "Configuration file")
	migrationsPath = flag.String("migrations", "./postgres/migrations", "Migrations path")
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
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		env.Postgres.User,
		env.Postgres.Password,
		env.Postgres.Host,
		env.Postgres.Port,
		env.Postgres.Database,
	)
}

func migratedb(dbpath string) error {
	db, err := sql.Open("postgres", dbpath)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	mpath := path.Clean(basepath + "./../../postgres/migrations")
	mpath = fmt.Sprintf("file://%v", mpath)
	m, err := migrate.NewWithDatabaseInstance(
		mpath, "postgres", driver,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	cfg := config.MustLoad(*env)
	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	url := parseDbUrl(cfg)
	err := migratedb(url)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	_, err = storage.New(url)
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

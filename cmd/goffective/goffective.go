package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"

	_ "github.com/OlegSadJktu/goffective/docs"
	"github.com/OlegSadJktu/goffective/internal/config"
	"github.com/OlegSadJktu/goffective/internal/dicontainer"
	localpg "github.com/OlegSadJktu/goffective/internal/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Goffective testing API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

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

func connectDb(url string) (*pg.DB, error) {
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
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
	db, err := connectDb(url)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	db.AddQueryHook(localpg.DBLogger{})
	container := dicontainer.New(db)
	defer container.Close()
	songsController := container.SongsController()

	router := gin.New()

	router.GET("/songs", songsController.Get)
	router.GET("/songs/:id", songsController.GetOne)
	router.POST("/songs", songsController.Create)
	router.DELETE("/songs/:id", songsController.Delete)
	router.PUT("/songs/:id", songsController.Update)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(fmt.Sprintf(":%v", cfg.Server.Port))
	slog.Error(err.Error())

}

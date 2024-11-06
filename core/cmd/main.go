package main

import (
	"encoding/gob"
	"gopher-toolbox/config"
	"log/slog"
	"net/http"
	"os"

	"gopher-toolbox/db"
	"gopher-toolbox/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/zepyrshut/rating-orama/internal/handlers"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

const version = "0.2.0-beta.20241116"

var app *config.App

func init() {
	gob.Register(map[string]string{})

	err := godotenv.Load()
	if err != nil {
		slog.Error("cannot load .env file", "error", err)
	}
	config.NewLogger(config.LogLevel(os.Getenv("LOG_LEVEL")))
	slog.Info("starting server")

	if os.Getenv("MIGRATE") == "true" {
		migrateDB()
	}
}

func migrateDB() {
	slog.Info("migrating database")
	m, err := migrate.New("file://database/migrations", os.Getenv("DATASOURCE"))
	if err != nil {
		slog.Error("cannot create migration", "error", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("cannot migrate", "error", err)
		panic(err)
	}
	if err == migrate.ErrNoChange {
		slog.Info("migration has no changes")
	}

	slog.Info("migration done")
}

func main() {

	app = &config.App{
		DataSource: os.Getenv("DATASOURCE"),
		UseCache:   utils.GetBool(os.Getenv("USE_CACHE")),
		AppInfo: config.AppInfo{
			GinMode: os.Getenv("GIN_MODE"),
			Version: version,
		},
	}

	dbPool := db.NewPostgresPool(app.DataSource)
	defer dbPool.Close()

	q := repository.NewPGXRepo(dbPool)
	h := handlers.New(q, app)
	r := Router(h, app)

	slog.Info("server started", "port", "8080", "version", version)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("cannot start server", "error", err)
	}
}

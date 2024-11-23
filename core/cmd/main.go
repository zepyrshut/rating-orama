package main

import (
	"embed"
	"encoding/gob"
	"gopher-toolbox/app"
	"log/slog"
	"ron"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zepyrshut/rating-orama/internal/handlers"
	"github.com/zepyrshut/rating-orama/internal/repository"
	"gopher-toolbox/db"
)

const version = "0.2.0-beta.20241116-4"

func init() {
	gob.Register(map[string]string{})
}

//go:embed database/migrations
var database embed.FS

func main() {
	app := app.New(version)
	app.Migrate(database)
	r := ron.New(func(e *ron.Engine) {
		e.Config.LogLevel = slog.LevelDebug
	})

	dbPool := db.NewPGXPool(app.Database.DataSource)
	defer dbPool.Close()

	q := repository.NewPGXRepo(dbPool)
	h := handlers.New(q, app)
	router(h, r)

	slog.Info("server started", "port", "8080", "version", version)
	err := r.Run(":8080")
	if err != nil {
		slog.Error("cannot start server", "error", err)
	}
}

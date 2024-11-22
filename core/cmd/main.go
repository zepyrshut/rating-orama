package main

import (
	"embed"
	"encoding/gob"
	"gopher-toolbox/app"
	"gopher-toolbox/db"
	"log/slog"
	"net/http"
	"ron"

	"github.com/zepyrshut/rating-orama/internal/handlers"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

//go:embed database/migrations
var database embed.FS

const version = "0.2.0-beta.20241116-4"

func init() {
	gob.Register(map[string]string{})
}

func main() {
	app := app.New(version)
	app.Migrate(database)
	r := ron.New(func(e *ron.Engine) {
		e.Render = ron.NewHTMLRender()
		e.Config.LogLevel = slog.LevelDebug
	})

	dbPool := db.NewPGXPool(app.Database.DataSource)
	defer dbPool.Close()

	q := repository.NewPGXRepo(dbPool)
	h := handlers.New(app, q)
	router(h, r)

	slog.Info("server started", "port", "8080", "version", version)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("cannot start server", "error", err)
	}
}

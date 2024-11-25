package main

import (
	"embed"
	"encoding/gob"
	"gopher-toolbox/app"
	"gopher-toolbox/db"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/zepyrshut/rating-orama/internal/handlers"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

//go:embed database/migrations
var database embed.FS

const version = "0.2.0-beta.20241116-4"
const appName = "rating-orama"

func init() {
	gob.Register(map[string]string{})
}

func main() {
	app := app.New(version)
	r := fiber.New(fiber.Config{
		AppName: appName,
	})
	
	dbPool := db.NewPGXPool(app.Database.DataSource)
	defer dbPool.Close()

	q := repository.NewPGXRepo(dbPool)
	h := handlers.New(app, q)
	router(h, r)

	slog.Info("server started", "port", "8080", "version", version)
	if err := r.Listen(":8080"); err != nil {
		slog.Error("cannot start server", "error", err)
	}
}

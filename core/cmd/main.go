package main

import (
	"embed"
	"encoding/gob"
	"gopher-toolbox/db"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/zepyrshut/rating-orama/internal/app"
	"github.com/zepyrshut/rating-orama/internal/handlers"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

const version = "0.2.3-beta.20250128-27"
const appName = "rating-orama"

func init() {
	gob.Register(map[string]string{})
}

//go:embed database/migrations
var database embed.FS

//go:embed templates
var templates embed.FS

func main() {
	engine := html.NewFileSystem(http.FS(templates), ".html")
	engine.Directory = "templates"

	app := app.NewExtendedApp(appName, version, ".env")
	app.Migrate(database)
	f := fiber.New(fiber.Config{
		AppName: appName,
		Views:   engine,
	})

	pgxPool := db.NewPGXPool(app.Database.DataSource)
	defer pgxPool.Close()

	r := repository.NewPGXRepo(pgxPool, app)
	h := handlers.New(r, app)
	router(h, f)

	slog.Info("server started", "port", "8080", "version", version)
	err := f.Listen(":8080")
	if err != nil {
		slog.Error("cannot start server", "error", err)
	}
}

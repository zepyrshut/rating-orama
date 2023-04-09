package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zepyrshut/rating-orama/handlers"
)

func Routes(app *fiber.App) {

	app.Use(recover.New())

	app.Static("/js", "./views/js")
	app.Static("/css", "./views/css")

	app.Get("/", handlers.Repo.Index)
	app.Get("/another", handlers.Repo.Another)
	app.Get("/tv-show", handlers.Repo.GetAllChapters)

	dev := app.Group("/dev")
	dev.Get("/ping", handlers.Repo.Ping)
	dev.Get("/panic", handlers.Repo.Panic)

}

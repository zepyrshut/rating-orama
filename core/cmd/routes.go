package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zepyrshut/rating-orama/internal/handlers"

	"gopher-toolbox/config"
)

func Router(h *handlers.Handlers, app *config.App) *gin.Engine {
	gin.SetMode(app.AppInfo.GinMode)
	r := gin.New()

	r.GET("/tvshow", h.GetTVShow)

	// app.Use(recover.New())

	// app.Static("/js", "./views/js")
	// app.Static("/css", "./views/css")

	// app.Get("/", handlers.Repo.Index)
	// app.Get("/another", handlers.Repo.Another)
	// app.Get("/tv-show", handlers.Repo.GetAllChapters)

	// dev := app.Group("/dev")
	// dev.Get("/ping", handlers.Repo.Ping)
	// dev.Get("/panic", handlers.Repo.Panic)

	return r
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zepyrshut/rating-orama/internal/handlers"

	"gopher-toolbox/config"
)

func Router(h *handlers.Handlers, app *config.App) *gin.Engine {
	gin.SetMode(app.AppInfo.GinMode)
	r := gin.New()

	r.GET("/ping", h.Ping)
	r.GET("/tvshow", h.GetTVShow)

	return r
}

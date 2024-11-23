package main

import (
	"github.com/zepyrshut/rating-orama/internal/handlers"
	"ron"
)

func router(h *handlers.Handlers, r *ron.Engine) {

	r.GET("/ping", h.Ping)
	r.GET("/error", h.Error)
	r.GET("/tvshow", h.GetTVShow)

}

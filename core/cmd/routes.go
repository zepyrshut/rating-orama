package main

import (
	"ron"

	"github.com/zepyrshut/rating-orama/internal/handlers"
)

func router(h *handlers.Handlers, r *ron.Engine) {

	r.GET("/ping", h.Ping)
	r.GET("/tvshow", h.GetTVShow)

}

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zepyrshut/rating-orama/internal/handlers"
)

func router(h *handlers.Handlers, r *fiber.App) {

	r.Get("/tvshow", h.GetTVShow)

}

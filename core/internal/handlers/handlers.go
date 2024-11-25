package handlers

import (
	"net/http"

	"gopher-toolbox/app"

	"github.com/gofiber/fiber/v2"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

type Handlers struct {
	app     *app.App
	queries repository.ExtendedQuerier
}

func New(app *app.App, q repository.ExtendedQuerier) *Handlers {
	return &Handlers{
		app:     app,
		queries: q,
	}
}

func (hq *Handlers) ToBeImplemented(c *fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).JSON("not implemented")
}

func (hq *Handlers) Ping(c *fiber.Ctx) error {
	return c.JSON("pong")
}

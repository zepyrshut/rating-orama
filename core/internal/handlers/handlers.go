package handlers

import (
	"context"
	"gopher-toolbox/app"
	"log/slog"
	"net/http"
	"ron"

	"github.com/zepyrshut/rating-orama/internal/repository"
)

type Handlers struct {
	App     *app.App
	Queries repository.ExtendedQuerier
}

func New(q repository.ExtendedQuerier, app *app.App) *Handlers {
	return &Handlers{
		Queries: q,
		App:     app,
	}
}

func (hq *Handlers) ToBeImplemented(c *ron.CTX, ctx context.Context) {
	c.JSON(http.StatusOK, ron.Data{
		"message": "Not implemented yet",
	})
}

func (hq *Handlers) Ping(c *ron.CTX, ctx context.Context) {
	slog.Info("ping", ron.RequestID, ctx.Value(ron.RequestID))
	c.JSON(http.StatusOK, ron.Data{
		"message": "pong",
	})
}

func (hq *Handlers) Error(c *ron.CTX, ctx context.Context) {
	slog.Error("error", ron.RequestID, ctx.Value(ron.RequestID))
	c.JSON(http.StatusInternalServerError, ron.Data{
		"req":     ctx.Value(ron.RequestID),
		"message": "error",
	})
}

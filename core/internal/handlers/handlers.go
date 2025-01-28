package handlers

import (
	"github.com/zepyrshut/rating-orama/internal/app"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

type Handlers struct {
	app     *app.ExtendedApp
	queries repository.ExtendedQuerier
}

func New(r repository.ExtendedQuerier, app *app.ExtendedApp) *Handlers {
	return &Handlers{
		app:     app,
		queries: r,
	}
}

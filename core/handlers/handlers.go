package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zepyrshut/rating-orama/app"
	"github.com/zepyrshut/rating-orama/repository"
)

type Repository struct {
	DB  repository.DBRepo
	App *app.Application
}

var Repo *Repository

func NewRepo(db *pgxpool.Pool, app *app.Application) *Repository {
	return &Repository{
		DB:  repository.NewPostgresRepo(db),
		App: app,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

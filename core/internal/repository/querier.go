package repository

import (
	"context"
	"github.com/zepyrshut/rating-orama/internal/scraper"

	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

type ExtendedQuerier interface {
	sqlc.Querier
	CreateTvShowWithEpisodes(ctx context.Context, tvShow sqlc.CreateTVShowParams, episodes []scraper.Episode) ([]sqlc.Episode, error)
}

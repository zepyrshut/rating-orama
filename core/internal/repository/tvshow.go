package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"

	"github.com/zepyrshut/rating-orama/internal/scraper"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

func (r *pgxRepository) CreateTvShowWithEpisodes(ctx context.Context, tvShow sqlc.CreateTVShowParams, episodes []scraper.Episode) ([]sqlc.Episode, error) {
	var sqlcEpisodes []sqlc.Episode
	err := r.execTx(ctx, func(tx pgx.Tx) error {
		qtx := r.WithTx(tx)
		tvShow, err := qtx.CreateTVShow(ctx, tvShow)
		if err != nil {
			return err
		}

		slog.Info("episodes lenght", "episodes", len(episodes))

		for _, episode := range episodes {
			sqlcEpisodeParams := episode.ToEpisodeParams(tvShow.ID)
			slog.Info("creating episode", "episode", sqlcEpisodeParams)
			episode, err := qtx.CreateEpisodes(ctx, sqlcEpisodeParams)
			if err != nil {
				return err
			}

			sqlcEpisodes = append(sqlcEpisodes, episode)
		}

		return nil
	})

	return sqlcEpisodes, err
}

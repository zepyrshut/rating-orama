package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
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

		for _, episode := range episodes {
			sqlcEpisodeParams := episode.ToEpisodeParams(tvShow.ID)
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

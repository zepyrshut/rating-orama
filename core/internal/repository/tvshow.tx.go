package repository

import (
	"context"

	"github.com/zepyrshut/rating-orama/internal/scraper"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

func (r *pgxRepository) CreateTvShowWithEpisodesTX(ctx context.Context, tvShow sqlc.CreateTVShowParams, episodes []scraper.Episode) ([]sqlc.Episode, error) {
	var err error
	var episodesSqlc []sqlc.Episode

	err = r.execTx(ctx, func(tx *sqlc.Queries) error {
		tvShow, err := tx.CreateTVShow(ctx, tvShow)
		if err != nil {
			return err
		}

		for _, episode := range episodes {
			sqlcEpisodeParams := episode.ToEpisodeParams(tvShow.ID)
			episode, err := tx.CreateEpisodes(ctx, sqlcEpisodeParams)
			if err != nil {
				return err
			}

			episodesSqlc = append(episodesSqlc, episode)
		}

		return nil
	})

	return episodesSqlc, err
}

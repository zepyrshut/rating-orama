package repository

import "github.com/zepyrshut/rating-orama/models"

type DBRepo interface {
	CheckIfTvShowExists(showID string) bool
	InsertTvShow(tvShow models.TvShow) error
	InsertEpisodes(tvShow models.TvShow) error
	FetchTvShow(showID string) (models.TvShow, error)
	IncreasePopularity(showID string)
	FetchEpisodes(showID string) ([]models.Season, error)
	TvShowAverageRating(show *models.TvShow)
	SeasonAverageRating(show *models.TvShow)
	TvShowMedianRating(show *models.TvShow)
	SeasonMedianRating(show *models.TvShow)
}

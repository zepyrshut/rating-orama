package repository

import (
	"context"
	"fmt"
	"github.com/zepyrshut/rating-orama/models"
	"time"
)

func (pg *postgresDBRepo) CheckIfTvShowExists(showID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT show_id FROM tv_show WHERE show_id = $1`

	var showIDFromDB string
	err := pg.DB.QueryRow(ctx, query, showID).Scan(&showIDFromDB)
	if err != nil {
		return false
	}

	return true
}

func (pg *postgresDBRepo) InsertTvShow(tvShow models.TvShow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryTvShow := `INSERT INTO tv_show (show_id, title, runtime) VALUES ($1, $2, $3)`

	_, err := pg.DB.Exec(ctx, queryTvShow, tvShow.ShowID, tvShow.Title, tvShow.Runtime)
	if err != nil {
		return err
	}

	err = pg.InsertEpisodes(tvShow)
	if err != nil {
		return err
	}

	return nil
}

func (pg *postgresDBRepo) InsertEpisodes(tvShow models.TvShow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `INSERT INTO episodes (episode_id, tv_show_id, season_number, title, number, aired, avg_rating, votes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for k, season := range tvShow.Seasons {
		for _, episode := range season.Episodes {
			_, err := pg.DB.Exec(ctx, query, episode.EpisodeID, tvShow.ShowID, k+1, episode.Title, episode.Number, episode.Aired, episode.AvgRating, episode.Votes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pg *postgresDBRepo) FetchTvShow(showID string) (models.TvShow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT show_id, title, runtime FROM tv_show WHERE show_id = $1`

	var tvShow models.TvShow
	var tvShowID int
	err := pg.DB.QueryRow(ctx, query, showID).Scan(&tvShowID, &tvShow.Title, &tvShow.Runtime)
	if err != nil {
		return tvShow, err
	}

	tvShow.ShowID = fmt.Sprintf("%07d", tvShowID)

	tvShow.Seasons, err = pg.FetchEpisodes(showID)
	if err != nil {
		return tvShow, err
	}

	pg.TvShowAverageRating(&tvShow)
	pg.SeasonAverageRating(&tvShow)

	pg.TvShowMedianRating(&tvShow)
	pg.SeasonMedianRating(&tvShow)

	pg.IncreasePopularity(showID)
	return tvShow, nil
}

func (pg *postgresDBRepo) IncreasePopularity(showID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `UPDATE tv_show SET popularity = popularity + 1 WHERE show_id = $1`

	_, err := pg.DB.Exec(ctx, query, showID)
	if err != nil {
		return
	}
}

func (pg *postgresDBRepo) FetchEpisodes(showID string) ([]models.Season, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT episode_id, season_number, title, number, aired, avg_rating, votes FROM episodes WHERE tv_show_id = $1 ORDER BY season_number, number`

	rows, err := pg.DB.Query(ctx, query, showID)
	if err != nil {
		return nil, err
	}

	var seasons []models.Season
	var episodeID int
	for rows.Next() {
		var episode models.Episode
		var seasonNumber int
		err = rows.Scan(&episodeID, &seasonNumber, &episode.Title, &episode.Number, &episode.Aired, &episode.AvgRating, &episode.Votes)
		if err != nil {
			return nil, err
		}

		episode.EpisodeID = fmt.Sprintf("%07d", episodeID)

		if len(seasons) < seasonNumber {
			seasons = append(seasons, models.Season{})
		}

		seasons[seasonNumber-1].Number = seasonNumber
		seasons[seasonNumber-1].Episodes = append(seasons[seasonNumber-1].Episodes, episode)
	}

	return seasons, nil
}

func (pg *postgresDBRepo) TvShowAverageRating(show *models.TvShow) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query :=
		`
SELECT 
    AVG(avg_rating), SUM(votes) 
FROM 
    episodes 
WHERE 
    tv_show_id = $1
AND
    votes > 0 AND avg_rating > 0
`

	var avgRating float64
	var votes int
	err := pg.DB.QueryRow(ctx, query, show.ShowID).Scan(&avgRating, &votes)
	if err != nil {
		return
	}

	show.AvgRating = avgRating
	show.Votes = votes
}

func (pg *postgresDBRepo) SeasonAverageRating(show *models.TvShow) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query :=
		`
SELECT 
    season_number, AVG(avg_rating), SUM(votes) 
FROM 
    episodes 
WHERE 
    tv_show_id = $1 
AND
	votes > 0 AND avg_rating > 0
GROUP BY 
    season_number 
ORDER BY 
    season_number;
`

	rows, err := pg.DB.Query(ctx, query, show.ShowID)
	if err != nil {
		return
	}

	for rows.Next() {
		var seasonNumber int
		var avgRating float64
		var votes int
		err = rows.Scan(&seasonNumber, &avgRating, &votes)
		if err != nil {
			return
		}

		show.Seasons[seasonNumber-1].AvgRating = avgRating
		show.Seasons[seasonNumber-1].Votes = votes
	}
}

func (pg *postgresDBRepo) TvShowMedianRating(show *models.TvShow) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query :=
		`
SELECT
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY avg_rating) AS median_rating
FROM
    episodes
WHERE
    tv_show_id = $1
AND
    votes > 0 AND avg_rating > 0;
`

	var medianRating float64
	err := pg.DB.QueryRow(ctx, query, show.ShowID).Scan(&medianRating)
	if err != nil {
		return
	}

	show.MedianRating = medianRating
}

func (pg *postgresDBRepo) SeasonMedianRating(show *models.TvShow) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query :=
		`
WITH episodes_with_ranks AS (
    SELECT
        episode_id,
        tv_show_id,
        season_number,
        avg_rating,
        votes,
        ROW_NUMBER() OVER (PARTITION BY season_number ORDER BY avg_rating) AS rank_asc,
        ROW_NUMBER() OVER (PARTITION BY season_number ORDER BY avg_rating DESC) AS rank_desc,
        COUNT(*) OVER (PARTITION BY season_number) AS season_episode_count
    FROM
        episodes
    WHERE
        tv_show_id = $1
	AND
		votes > 0 AND avg_rating > 0
),
episodes_filtered AS (
    SELECT
        episode_id,
        tv_show_id,
        season_number,
        avg_rating,
        votes
    FROM
        episodes_with_ranks
    WHERE
        rank_asc > season_episode_count * 0.01 AND rank_desc > season_episode_count * 0.01
),
seasons AS (
    SELECT DISTINCT
        season_number
    FROM
        episodes_filtered
)
SELECT
    seasons.season_number,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY episodes_filtered.avg_rating) AS median_rating
FROM
    seasons
JOIN
    episodes_filtered ON seasons.season_number = episodes_filtered.season_number
GROUP BY
    seasons.season_number
ORDER BY
    seasons.season_number;
`

	rows, err := pg.DB.Query(ctx, query, show.ShowID)
	if err != nil {
		return
	}

	for rows.Next() {
		var seasonNumber int
		var medianRating float64
		err = rows.Scan(&seasonNumber, &medianRating)
		if err != nil {
			return
		}

		show.Seasons[seasonNumber-1].MedianRating = medianRating
	}

}

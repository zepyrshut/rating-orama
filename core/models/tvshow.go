package models

import (
	"strconv"
	"time"
)

type Popularity struct {
	ShowID      string `json:"show_id"`
	TimesViewed int    `json:"times_viewed"`
}

type TvShow struct {
	ShowID       string   `json:"show_id"`
	Title        string   `json:"title"`
	Runtime      int      `json:"runtime"`
	Votes        int      `json:"votes"`
	AvgRating    float64  `json:"avg_rating"`
	MedianRating float64  `json:"median_rating"`
	Seasons      []Season `json:"seasons"`
}

type Season struct {
	Number       int       `json:"number"`
	AvgRating    float64   `json:"avg_rating"`
	MedianRating float64   `json:"median_rating"`
	Votes        int       `json:"votes"`
	Episodes     []Episode `json:"episodes"`
}

type Episode struct {
	Number    int       `json:"number"`
	EpisodeID string    `json:"episode_id"`
	Title     string    `json:"title"`
	Aired     time.Time `json:"aired"`
	AvgRating float64   `json:"avg_rating"`
	Votes     int       `json:"votes"`
}

func (tvShow *TvShow) TvShowBuilder(tvShowDTO TvShowDTO) {
	tvShow.ShowID = tvShowDTO.ShowID
	tvShow.Title = tvShowDTO.Title
	tvShow.Runtime, _ = strconv.Atoi(tvShowDTO.Runtime)

	lastSeasonNumber := tvShowDTO.Episodes[len(tvShowDTO.Episodes)-1].SeasonID
	if lastSeasonNumber == -1 {
		lastSeasonNumber = tvShowDTO.Episodes[len(tvShowDTO.Episodes)-2].SeasonID
	}
	seasons := make([]Season, lastSeasonNumber)

	for currentSeason := 1; currentSeason <= lastSeasonNumber; currentSeason++ {
		for _, episode := range tvShowDTO.Episodes {
			if episode.SeasonID == currentSeason {
				seasons[currentSeason-1].Number = currentSeason
				seasons[currentSeason-1].Episodes = append(seasons[currentSeason-1].Episodes, Episode{
					Number:    episode.Number,
					EpisodeID: episode.EpisodeID,
					Title:     episode.Title,
					Aired:     episode.Aired.Time,
					AvgRating: episode.AvgRating,
					Votes:     episode.Votes,
				})
			}
		}
	}

	tvShow.Seasons = seasons
}

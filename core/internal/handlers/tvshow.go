package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/zepyrshut/rating-orama/internal/scraper"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

func (hq *Handlers) GetIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Rating Orama",
	}, "layouts/main")
}

func (hq *Handlers) GetTVShow(c *fiber.Ctx) error {
	ttShowID := c.Query("ttid")

	if ttShowID == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	var title string
	var scraperEpisodes []scraper.Episode
	var sqlcEpisodes []sqlc.Episode
	var totalVoteCount int32

	sqlcTvShow, err := hq.queries.CheckTVShowExists(c.Context(), ttShowID)
	if err != nil {
		title, scraperEpisodes = scraper.ScrapeEpisodes(ttShowID)
		//TODO: make transactional
		ttShow, err := hq.queries.CreateTVShow(c.Context(), sqlc.CreateTVShowParams{
			TtImdb: ttShowID,
			Name:   title,
		})
		if err != nil {
			slog.Error("failed to create tv show", "ttid", ttShowID, "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		slog.Info("ttshowid", "id", ttShow.ID)
		for _, episode := range scraperEpisodes {
			sqlcEpisodesParams := episode.ToEpisodeParams(ttShow.ID)
			sqlcEpisode, err := hq.queries.CreateEpisodes(c.Context(), sqlcEpisodesParams)
			if err != nil {
				slog.Error("failed to create episodes", "ttid", ttShowID, "error", err)
				return c.SendStatus(http.StatusInternalServerError)
			}
			sqlcEpisodes = append(sqlcEpisodes, sqlcEpisode)
		}

		slog.Info("scraped seasons", "ttid", ttShowID, "title", title)
	} else {
		title = sqlcTvShow.Name
		sqlcEpisodes, err = hq.queries.GetEpisodes(c.Context(), sqlcTvShow.ID)
		if err != nil {
			slog.Error("failed to get episodes", "ttid", ttShowID, "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		for _, episode := range sqlcEpisodes {
			totalVoteCount += episode.VoteCount
		}

		hq.queries.IncreasePopularity(c.Context(), ttShowID)
		slog.Info("tv show exists", "ttid", ttShowID, "title", sqlcTvShow.Name)
	}

	episodesJSON, err := json.Marshal(sqlcEpisodes)
	if err != nil {
		slog.Error("failed to marshal episodes", "ttid", ttShowID, "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	// calculate avg rating for the show
	avgRatingShow, err := hq.queries.TvShowAverageRating(c.Context(), sqlcTvShow.ID)
	if err != nil {
		slog.Error("failed to calculate avg rating for the show", "ttid", ttShowID, "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	medianRatingShow, err := hq.queries.TvShowMedianRating(c.Context(), sqlcTvShow.ID)
	if err != nil {
		slog.Error("failed to calculate median rating for the show", "ttid", ttShowID, "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	seasongAvgRatings, err := hq.queries.SeasonAverageRating(c.Context(), sqlc.SeasonAverageRatingParams{
		TvShowID: sqlcTvShow.ID,
	})

	seasonMedianRatings, err := hq.queries.SeasonMedianRating(c.Context(), sqlc.SeasonMedianRatingParams{
		TvShowID: sqlcTvShow.ID,
	})

	return c.Render("tvshow", fiber.Map{
		"Title":                 sqlcTvShow.Name,
		"tvshow":                sqlcTvShow,
		"episodes":              string(episodesJSON),
		"avg_rating_show":       avgRatingShow,
		"median_rating_show":    medianRatingShow,
		"season_avg_ratings":    seasongAvgRatings,
		"season_median_ratings": seasonMedianRatings,
		"total_vote_count":      totalVoteCount,
	}, "layouts/main")
}

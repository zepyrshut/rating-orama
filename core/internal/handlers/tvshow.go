package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/zepyrshut/rating-orama/internal/scraper"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

func (hq *Handlers) GetTVShow(c *fiber.Ctx) error {
	ttShowID := c.Query("ttid")

	var title string
	var scraperEpisodes []scraper.Episode
	var sqlcEpisodes []sqlc.Episode

	tvShow, err := hq.queries.CheckTVShowExists(c.Context(), ttShowID)
	if err != nil {
		title, scraperEpisodes = scraper.ScrapeEpisodes(ttShowID)
		// TODO: make transactional
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
		title = tvShow.Name
		sqlcEpisodes, err = hq.queries.GetEpisodes(c.Context(), tvShow.ID)
		if err != nil {
			slog.Error("failed to get episodes", "ttid", ttShowID, "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		hq.queries.IncreasePopularity(c.Context(), ttShowID)
		slog.Info("tv show exists", "ttid", ttShowID, "title", tvShow.Name)
	}

	return c.JSON(fiber.Map{
		"popularity": tvShow.Popularity,
		"title":      title,
		"seasons":    sqlcEpisodes,
	})
}

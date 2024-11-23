package handlers

import (
	"context"
	"gopher-toolbox/app"
	"log/slog"
	"net/http"
	"ron"

	"github.com/zepyrshut/rating-orama/internal/scraper"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

func (hq *Handlers) GetTVShow(c *ron.CTX, ctx context.Context) {
	ttShowID := c.Query("ttid")
	slog.Info("", "ttid", ttShowID, ron.RequestID, ctx.Value(ron.RequestID))

	var title string
	var scraperEpisodes []scraper.Episode
	var sqlcEpisodes []sqlc.Episode

	tvShow, err := hq.Queries.CheckTVShowExists(ctx, ttShowID)
	if err != nil {
		title, scraperEpisodes = scraper.ScrapeEpisodes(ttShowID)

		sqlcEpisodes, err = hq.Queries.CreateTvShowWithEpisodes(ctx, sqlc.CreateTVShowParams{
			TtImdb: ttShowID,
			Name:   title,
		}, scraperEpisodes)
		if err != nil {
			slog.Error("failed to create tv show with episodes", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, ron.Data{"error": app.ErrorCreating})
			return
		}

		slog.Info("scraped seasons", "ttid", ttShowID, "title", title)
	} else {
		title = tvShow.Name
		sqlcEpisodes, err = hq.Queries.GetEpisodes(ctx, tvShow.ID)
		if err != nil {
			slog.Error("failed to get episodes", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, ron.Data{"error": app.ErrorGetting})
			return
		}

		if err := hq.Queries.IncreasePopularity(ctx, ttShowID); err != nil {
			slog.Error("failed to increase popularity", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, ron.Data{"error": app.ErrorUpdating})
			return
		}

		slog.Info("tv show exists", "ttid", ttShowID, "title", tvShow.Name)
	}

	tvShowMedian, _ := hq.Queries.TvShowMedianRating(ctx, sqlcEpisodes[0].TvShowID)
	tvShowAverage, _ := hq.Queries.TvShowAverageRating(ctx, sqlcEpisodes[0].TvShowID)

	c.JSON(http.StatusOK, ron.Data{
		"popularity":    tvShow.Popularity,
		"title":         title,
		"seasons":       sqlcEpisodes,
		"tvShowMedian":  tvShowMedian,
		"tvShowAverage": tvShowAverage,
	})
}

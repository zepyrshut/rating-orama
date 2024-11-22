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
	//slog.Info("", "ttid", ttShowID, RequestID, ctx.Value(RequestID))

	var title string
	var scraperEpisodes []scraper.Episode
	var sqlcEpisodes []sqlc.Episode

	tvShow, err := hq.queries.CheckTVShowExists(ctx, ttShowID)
	if err != nil {
		title, scraperEpisodes = scraper.ScrapeEpisodes(ttShowID)
<<<<<<< Updated upstream

		sqlcEpisodes, err = hq.Queries.CreateTvShowWithEpisodes(c, sqlc.CreateTVShowParams{
=======
		// TODO: make transactional
		ttShow, err := hq.queries.CreateTVShow(ctx, sqlc.CreateTVShowParams{
>>>>>>> Stashed changes
			TtImdb: ttShowID,
			Name:   title,
		}, scraperEpisodes)
		if err != nil {
<<<<<<< Updated upstream
			slog.Error("failed to create tv show with episodes", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorCreating})
			return
=======
			slog.Error("failed to create tv show", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, ron.Data{"error": app.ErrorCreating})
		}

		slog.Info("ttshowid", "id", ttShow.ID)
		for _, episode := range scraperEpisodes {
			sqlcEpisodesParams := episode.ToEpisodeParams(ttShow.ID)
			sqlcEpisode, err := hq.queries.CreateEpisodes(ctx, sqlcEpisodesParams)
			if err != nil {
				slog.Error("failed to create episodes", "ttid", ttShowID, "error", err)
				c.JSON(http.StatusInternalServerError, ron.Data{"error": app.ErrorCreating})
				return
			}

			sqlcEpisodes = append(sqlcEpisodes, sqlcEpisode)
>>>>>>> Stashed changes
		}

		slog.Info("scraped seasons", "ttid", ttShowID, "title", title)
	} else {
		title = tvShow.Name
		sqlcEpisodes, err = hq.queries.GetEpisodes(ctx, tvShow.ID)
		if err != nil {
			slog.Error("failed to get episodes", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, ron.Data{"error": app.ErrorGetting})
			return
		}

<<<<<<< Updated upstream
		if err := hq.Queries.IncreasePopularity(c, ttShowID); err != nil {
			slog.Error("failed to increase popularity", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorUpdating})
			return
		}

		slog.Info("tv show exists", "ttid", ttShowID, "title", tvShow.Name)
	}

	tvShowMedian, _ := hq.Queries.TvShowMedianRating(c, sqlcEpisodes[0].TvShowID)
	tvShowAverage, _ := hq.Queries.TvShowAverageRating(c, sqlcEpisodes[0].TvShowID)

	c.JSON(http.StatusOK, gin.H{
		"popularity":    tvShow.Popularity,
		"title":         title,
		"seasons":       sqlcEpisodes,
		"tvShowMedian":  tvShowMedian,
		"tvShowAverage": tvShowAverage,
=======
		hq.queries.IncreasePopularity(ctx, ttShowID)
		slog.Info("tv show exists", "ttid", ttShowID, "title", tvShow.Name)
	}

	c.JSON(http.StatusOK, ron.Data{
		"popularity": tvShow.Popularity,
		"title":      title,
		"seasons":    sqlcEpisodes,
>>>>>>> Stashed changes
	})
}

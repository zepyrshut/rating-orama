package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zepyrshut/rating-orama/internal/scraper"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

func (hq *Handlers) GetTVShow(c *gin.Context) {
	ttShowID := c.Query("ttid")
	slog.Info("", "ttid", ttShowID, RequestID, c.Request.Context().Value(RequestID))

	var title string
	var scraperEpisodes []scraper.Episode
	var sqlcEpisodes []sqlc.Episode

	tvShow, err := hq.Queries.CheckTVShowExists(c, ttShowID)
	if err != nil {
		title, scraperEpisodes = scraper.ScrapeEpisodes(ttShowID)

		sqlcEpisodes, err = hq.Queries.CreateTvShowWithEpisodes(c, sqlc.CreateTVShowParams{
			TtImdb: ttShowID,
			Name:   title,
		}, scraperEpisodes)
		if err != nil {
			slog.Error("failed to create tv show with episodes", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorCreating})
			return
		}

		slog.Info("scraped seasons", "ttid", ttShowID, "title", title)
	} else {
		title = tvShow.Name
		sqlcEpisodes, err = hq.Queries.GetEpisodes(c, tvShow.ID)
		if err != nil {
			slog.Error("failed to get episodes", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorGetting})
			return
		}

		if err := hq.Queries.IncreasePopularity(c, ttShowID); err != nil {
			slog.Error("failed to increase popularity", "ttid", ttShowID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrorUpdating})
			return
		}

		slog.Info("tv show exists", "ttid", ttShowID, "title", tvShow.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"popularity": tvShow.Popularity,
		"title":      title,
		"seasons":    sqlcEpisodes,
	})
}

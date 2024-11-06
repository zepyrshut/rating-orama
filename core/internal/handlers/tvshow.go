package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zepyrshut/rating-orama/internal/scraper"
)

func (hq *Handlers) GetTVShow(c *gin.Context) {
	ttShowID := c.Query("ttid")
	slog.Info("GetTVShow", "ttid", ttShowID)

	title, seasons := scraper.ScrapeSeasons(ttShowID)

	slog.Info("scraped seasons", "ttid", ttShowID, "title", title)

	c.JSON(http.StatusOK, gin.H{
		"title":   title,
		"seasons": seasons,
	})
}

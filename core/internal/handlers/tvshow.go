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

	seasons := scraper.ScrapeSeasons(ttShowID)

	c.JSON(http.StatusOK, gin.H{
		"seasons": seasons,
	})
}

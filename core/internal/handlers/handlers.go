package handlers

import (
	"context"
	"net/http"
	"ron"

	"gopher-toolbox/app"

	"github.com/zepyrshut/rating-orama/internal/repository"
)

type Handlers struct {
	app     *app.App
	queries repository.ExtendedQuerier
}

func New(app *app.App, q repository.ExtendedQuerier) *Handlers {
	return &Handlers{
		app:     app,
		queries: q,
	}
}

func (hq *Handlers) ToBeImplemented(c *ron.CTX, ctx context.Context) {
	c.JSON(http.StatusOK, ron.Data{
		"message": "To be implemented",
	})
}

<<<<<<< Updated upstream
func (hq *Handlers) Ping(c *gin.Context) {
	slog.Info("ping", RequestID, c.Request.Context().Value(RequestID))
	c.JSON(http.StatusOK, gin.H{
=======
func (hq *Handlers) Ping(c *ron.CTX, ctx context.Context) {
	c.JSON(http.StatusOK, ron.Data{
>>>>>>> Stashed changes
		"message": "pong",
	})
}

// // TODO: Extract to toolbox
// func handleQueryError(c *gin.Context, err error, errorMap map[string]string, logMessage string, defaultErrorMessage string) bool {
// 	if err != nil {
// 		for key, message := range errorMap {
// 			if strings.Contains(err.Error(), key) {
// 				slog.Error(logMessage, "error", message, RequestID, c.Request.Context().Value(RequestID))
// 				c.JSON(http.StatusConflict, gin.H{"error": message})
// 				return true
// 			}
// 		}

// 		slog.Error(logMessage, "error", err.Error(), RequestID, c.Request.Context().Value(RequestID))
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrorMessage})
// 		return true
// 	}
// 	return false
// }

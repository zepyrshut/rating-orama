package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"gopher-toolbox/config"

	"github.com/gin-gonic/gin"
	"github.com/zepyrshut/rating-orama/internal/repository"
)

// TODO: Extract to toolbox
const (
	InvalidRequest  string = "invalid_request"
	InternalError   string = "internal_error"
	RequestID       string = "request_id"
	NotFound        string = "not_found"
	Created         string = "created"
	Updated         string = "updated"
	Deleted         string = "deleted"
	Enabled         string = "enabled"
	Disabled        string = "disabled"
	Retrieved       string = "retrieved"
	ErrorCreating   string = "error_creating"
	ErrorUpdating   string = "error_updating"
	ErrorEnabling   string = "error_enabling"
	ErrorDisabling  string = "error_disabling"
	ErrorGetting    string = "error_getting"
	ErrorGettingAll string = "error_getting_all"
	InvalidEntityID string = "invalid_entity_id"
	NotImplemented  string = "not_implemented"

	UserUsernameKey       string = "user_username_key"
	UserEmailKey          string = "user_email_key"
	UsernameAlReadyExists string = "username_already_exists"
	EmailAlreadyExists    string = "email_already_exists"
	IncorrectPassword     string = "incorrect_password"
	ErrorGeneratingToken  string = "error_generating_token"
	LoggedIn              string = "logged_in"

	CategoryNameKey       string = "category_name_key"
	CategoryAlreadyExists string = "category_already_exists"

	ItemsNameKey      string = "items_name_key"
	NameAlreadyExists string = "name_already_exists"
)

type Handlers struct {
	App     *config.App
	Queries repository.ExtendedQuerier
}

func New(q repository.ExtendedQuerier, app *config.App) *Handlers {
	return &Handlers{
		Queries: q,
		App:     app,
	}
}

func (hq *Handlers) ToBeImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}

func (hq *Handlers) Ping(c *gin.Context) {
	slog.Info("ping", RequestID, c.Request.Context().Value(RequestID))
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// TODO: Extract to toolbox
func handleQueryError(c *gin.Context, err error, errorMap map[string]string, logMessage string, defaultErrorMessage string) bool {
	if err != nil {
		for key, message := range errorMap {
			if strings.Contains(err.Error(), key) {
				slog.Error(logMessage, "error", message, RequestID, c.Request.Context().Value(RequestID))
				c.JSON(http.StatusConflict, gin.H{"error": message})
				return true
			}
		}

		slog.Error(logMessage, "error", err.Error(), RequestID, c.Request.Context().Value(RequestID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrorMessage})
		return true
	}
	return false
}

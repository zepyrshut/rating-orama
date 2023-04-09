package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zepyrshut/rating-orama/models"
	"io"
	"net/http"
)

func (rp Repository) GetAllChapters(c *fiber.Ctx) error {
	tvShow := models.TvShow{}

	ttShowID := c.Query("id")

	if ttShowID[0:2] == "tt" {
		ttShowID = ttShowID[2:]
	}

	exist := rp.DB.CheckIfTvShowExists(ttShowID)

	if !exist {
		url := fmt.Sprintf(rp.App.Environment.HarvesterApi, ttShowID)
		response, _ := http.Get(url)
		body, _ := io.ReadAll(response.Body)
		err := json.Unmarshal(body, &tvShow)
		if err != nil {
			rp.App.Error(err.Error())
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
		err = rp.DB.InsertTvShow(tvShow)
		if err != nil {
			rp.App.Error(err.Error())
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
	}

	tvShow, err := rp.DB.FetchTvShow(ttShowID)
	if err != nil {
		rp.App.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	tvShowJSON, err := json.Marshal(tvShow)
	if err != nil {
		rp.App.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Render("charts", fiber.Map{
		"TvShow":     tvShow,
		"TvShowJSON": string(tvShowJSON),
	})
}

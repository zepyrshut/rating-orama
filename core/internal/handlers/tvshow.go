package handlers

//func (hq *Handlers) GetAllChapters(c *fiber.Ctx) error {
// tvShow := models.TvShow{}

// ttShowID := c.Query("id")

// if ttShowID[0:2] == "tt" {
// 	ttShowID = ttShowID[2:]
// }

// exist := hq.DB.CheckIfTvShowExists(ttShowID)

// if !exist {
// 	url := fmt.Sprintf(hq.App.Environment.HarvesterApi, ttShowID)
// 	response, _ := http.Get(url)
// 	body, _ := io.ReadAll(response.Body)
// 	err := json.Unmarshal(body, &tvShow)
// 	if err != nil {
// 		hq.App.Error(err.Error())
// 		return c.Status(http.StatusInternalServerError).JSON(err)
// 	}
// 	err = hq.DB.InsertTvShow(tvShow)
// 	if err != nil {
// 		hq.App.Error(err.Error())
// 		return c.Status(http.StatusInternalServerError).JSON(err)
// 	}
// }

// tvShow, err := hq.DB.FetchTvShow(ttShowID)
// if err != nil {
// 	hq.App.Error(err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(err)
// }

// tvShowJSON, err := json.Marshal(tvShow)
// if err != nil {
// 	hq.App.Error(err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(err)
// }

// return c.Render("charts", fiber.Map{
// 	"TvShow":     tvShow,
// 	"TvShowJSON": string(tvShowJSON),
// })
//}

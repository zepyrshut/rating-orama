package models

// type TvShowDTO struct {
// 	ShowID   string       `json:"tt_show_id"`
// 	Title    string       `json:"title"`
// 	Runtime  string       `json:"runtime"`
// 	Episodes []EpisodeDTO `json:"episodes"`
// }

// type EpisodeDTO struct {
// 	Number    int       `json:"number"`
// 	SeasonID  int       `json:"season_id"`
// 	EpisodeID string    `json:"tt_episode_id"`
// 	Title     string    `json:"title"`
// 	Aired     AiredTime `json:"aired"`
// 	AvgRating float64   `json:"avg_rating"`
// 	Votes     int       `json:"votes"`
// }

// type AiredTime struct {
// 	time.Time
// }

// func (tvShow *TvShow) UnmarshalJSON(data []byte) error {
// 	var tvShowDTO TvShowDTO
// 	err := json.Unmarshal(data, &tvShowDTO)
// 	if err != nil {
// 		return err
// 	}

// 	tvShow.TvShowBuilder(tvShowDTO)
// 	return nil
// }

// func (aired *AiredTime) UnmarshalJSON(data []byte) error {
// 	if string(data) == "null" || string(data) == "" {
// 		return nil
// 	}

// 	var s string
// 	if err := json.Unmarshal(data, &s); err != nil {
// 		return nil
// 	}

// 	t, err := utils.TimeParser(s)
// 	if err != nil {
// 		return err
// 	}

// 	aired.Time = t
// 	return nil
// }

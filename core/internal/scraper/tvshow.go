package scraper

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zepyrshut/rating-orama/internal/sqlc"
)

type Episode struct {
	Season    int
	Episode   int
	Released  time.Time
	Name      string
	Plot      string
	Rate      float32
	VoteCount int
}

func (e Episode) ToEpisodeParams(tvShowID int32) sqlc.CreateEpisodesParams {

	var date pgtype.Date
	_ = date.Scan(e.Released)

	return sqlc.CreateEpisodesParams{
		TvShowID:  tvShowID,
		Season:    int32(e.Season),
		Episode:   int32(e.Episode),
		Name:      e.Name,
		Released:  date,
		Plot:      e.Plot,
		AvgRating: e.Rate,
		VoteCount: int32(e.VoteCount),
	}
}

func ScrapeEpisodes(ttImdb string) (string, []Episode) {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US")
	})

	var allSeasons []Episode
	var seasons []int
	var title string

	c.OnHTML(os.Getenv("SEASON_SELECTOR"), func(e *colly.HTMLElement) {
		seasonText := strings.TrimSpace(e.Text)
		seasonNum, err := strconv.Atoi(seasonText)
		if err == nil {
			seasons = append(seasons, seasonNum)
		}
	})

	c.OnHTML(os.Getenv("TITLE_SELECTOR"), func(e *colly.HTMLElement) {
		title = e.Text
	})

	c.OnScraped(func(r *colly.Response) {
		seasonMap := make(map[int]bool)
		var uniqueSeasons []int
		slog.Info("scraped seasons", "seasons", seasons)
		for _, seasonNum := range seasons {
			if !seasonMap[seasonNum] {
				seasonMap[seasonNum] = true
				uniqueSeasons = append(uniqueSeasons, seasonNum)
			}
		}

		sort.Ints(uniqueSeasons)
		episodeCollector := c.Clone()

		episodeCollector.OnResponse(func(r *colly.Response) {
			slog.Info("response", "url", r.Request.URL)
			season := extractEpisodesFromSeason(string(r.Body))
			allSeasons = append(allSeasons, season...)
		})

		for _, seasonNum := range uniqueSeasons {
			seasonURL := fmt.Sprintf(os.Getenv("IMDB_EPISODES_URL"), ttImdb, seasonNum)
			slog.Info("visiting season", "url", seasonURL)
			_ = episodeCollector.Visit(seasonURL)
		}

		episodeCollector.Wait()
	})

	_ = c.Visit(fmt.Sprintf(os.Getenv("VISIT_URL"), ttImdb))
	c.Wait()

	slog.Info("scraped all seasons", "length", len(allSeasons))
	return title, allSeasons
}

func extractEpisodesFromSeason(data string) []Episode {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		slog.Error("error parsing html")
		return []Episode{}
	}

	var episodes []Episode
	doc.Find(os.Getenv("EPISODE_CARD_SELECTOR")).Each(func(i int, s *goquery.Selection) {
		var episode Episode

		seasonEpisodeTitle := s.Find(os.Getenv("SEASON_EPISODE_AND_TITLE_SELECTOR")).Text()
		episode.Season, episode.Episode, episode.Name = parseSeasonEpisodeTitle(seasonEpisodeTitle)

		releasedDate := s.Find(os.Getenv("RELEASED_DATE_SELECTOR")).Text()
		episode.Released = parseReleasedDate(releasedDate)

		plot := s.Find(os.Getenv("PLOT_SELECTOR")).Text()
		if plot == "Add a plot" {
			episode.Plot = ""
		} else {
			episode.Plot = plot
		}

		starRating := s.Find(os.Getenv("STAR_RATING_SELECTOR")).Text()
		episode.Rate = parseStarRating(starRating)

		voteCount := s.Find(os.Getenv("VOTE_COUNT_SELECTOR")).Text()
		episode.VoteCount = parseVoteCount(voteCount)

		episodes = append(episodes, episode)
	})

	slog.Info("extracted episodes", "length", len(episodes))
	return episodes
}

func parseSeasonEpisodeTitle(input string) (int, int, string) {
	re := regexp.MustCompile(`S(\d+)\.E(\d+)\s*∙\s*(.+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 4 {
		return 0, 0, ""
	}

	seasonNum, err1 := strconv.Atoi(matches[1])
	episodeNum, err2 := strconv.Atoi(matches[2])
	name := strings.TrimSpace(matches[3])

	if err1 != nil || err2 != nil {
		return 0, 0, ""
	}

	return seasonNum, episodeNum, name
}

func parseReleasedDate(releasedDate string) time.Time {
	const layout = "Mon, Jan 2, 2006"
	parsedDate, err := time.Parse(layout, releasedDate)
	if err != nil {
		slog.Error("error parsing date", "date", releasedDate)
		return time.Time{}
	}
	return parsedDate
}

func parseStarRating(starRating string) float32 {
	rating, err := strconv.ParseFloat(starRating, 32)
	if err != nil || rating < 0 || rating > 10 {
		slog.Warn("error parsing rating, out of limits", "rating", starRating)
		return 0
	}
	return float32(rating)
}

func parseVoteCount(voteCount string) int {
	re := regexp.MustCompile(`\(([\d.]+)(K?)\)`)
	matches := re.FindStringSubmatch(voteCount)
	if len(matches) != 3 {
		slog.Error("error parsing vote count", "count", voteCount)
		return 0
	}

	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		slog.Error("error parsing vote count", "count", voteCount)
		return 0
	}

	if matches[2] == "K" {
		num *= 1000
	}

	return int(num)
}

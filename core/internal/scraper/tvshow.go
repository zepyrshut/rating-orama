package scraper

import (
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Episode struct {
	Season    int
	Episode   int
	Released  time.Time
	Name      string
	Plot      string
	Rate      float64
	VoteCount int
}

type Season []Episode

const seasonsSelector = "ul.ipc-tabs a[data-testid='tab-season-entry']"
const episodesSelector = "section.sc-1e7f96be-0.ZaQIL"
const nextSeasonButtonSelector = "#next-season-btn"
const imdbEpisodesURL = "https://www.imdb.com/title/%s/episodes?season=%d"

func ScrapeSeasons(ttImdb string) []Season {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	var allSeasons []Season
	var seasons []int

	c.OnHTML("ul.ipc-tabs a[data-testid='tab-season-entry']", func(e *colly.HTMLElement) {
		seasonText := strings.TrimSpace(e.Text)
		seasonNum, err := strconv.Atoi(seasonText)
		if err == nil {
			seasons = append(seasons, seasonNum)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		seasonMap := make(map[int]bool)
		uniqueSeasons := []int{}
		slog.Info("scraped seasons", "seasons", seasons)
		for _, seasonNum := range seasons {
			if !seasonMap[seasonNum] {
				seasonMap[seasonNum] = true
				uniqueSeasons = append(uniqueSeasons, seasonNum)
			}
		}

		sort.Ints(uniqueSeasons)

		episodeCollector := c.Clone()

		episodeCollector.OnHTML(episodesSelector, func(e *colly.HTMLElement) {
			seasonEpisodes := extractEpisodesFromSeason(e.Text)
			allSeasons = append(allSeasons, seasonEpisodes)
			//allEpisodes = append(allEpisodes, seasonEpisodes...)
		})

		for _, seasonNum := range uniqueSeasons {
			seasonURL := fmt.Sprintf(imdbEpisodesURL, ttImdb, seasonNum)
			slog.Info("visiting season", "url", seasonURL)
			episodeCollector.Visit(seasonURL)
		}

		episodeCollector.Wait()
	})

	c.Visit(fmt.Sprintf("https://www.imdb.com/title/%s/episodes", ttImdb))
	c.Wait()

	return allSeasons
}

func extractEpisodesFromSeason(data string) Season {
	const pattern = `(S\d+\.E\d+)\s∙\s(.*?)` +
		`(Mon|Tue|Wed|Thu|Fri|Sat|Sun),\s` +
		`(.*?),\s(\d{4})(.*?)` +
		`(\d\.\d{1,2}\/10) \((\d+K)\)Rate`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(data, -1)

	episodes := make([]Episode, 0, len(matches))

	for _, match := range matches {
		var episode Episode

		seasonEpisode := match[1]

		name := strings.TrimSpace(match[2])

		day := match[3]
		dateRest := strings.TrimSpace(match[4])
		year := match[5]

		plot := strings.TrimSpace(match[6])
		rate := match[7]
		voteCount := match[8]

		seasonNum := strings.TrimPrefix(strings.Split(seasonEpisode, ".")[0], "S")
		episodeNum := strings.TrimPrefix(strings.Split(seasonEpisode, ".")[1], "E")

		votes, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(voteCount, "K"), "K"))

		episode.Name = name
		episode.Episode, _ = strconv.Atoi(episodeNum)
		episode.Season, _ = strconv.Atoi(seasonNum)
		episode.Released, _ = time.Parse("Mon, Jan 2, 2006", fmt.Sprintf("%s, %s, %s", day, dateRest, year))
		episode.Plot = plot
		episode.Rate, _ = strconv.ParseFloat(strings.TrimSuffix(rate, "/10"), 2)
		episode.VoteCount = votes * 1000

		episodes = append(episodes, episode)
	}

	return episodes
}

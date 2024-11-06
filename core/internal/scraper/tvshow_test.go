package scraper

import (
	"testing"
	"time"
)

func Test_parseSeasonEpisodeTitle(t *testing.T) {
	var tests = []struct {
		given    string
		expected struct {
			seasonNum  int
			episodeNum int
			name       string
		}
	}{
		{"S5.E1 ∙ Live Free or Die", struct {
			seasonNum  int
			episodeNum int
			name       string
		}{5, 1, "Live Free or Die"}},
		{"S5.E13 ∙ To'hajiilee", struct {
			seasonNum  int
			episodeNum int
			name       string
		}{5, 13, "To'hajiilee"}},
	}

	for _, tt := range tests {
		t.Run(tt.given, func(t *testing.T) {
			seasonNum, episodeNum, name := parseSeasonEpisodeTitle(tt.given)
			if seasonNum != tt.expected.seasonNum || episodeNum != tt.expected.episodeNum || name != tt.expected.name {
				t.Errorf("parseSeasonEpisodeTitle(%s): expected %d, %d, %s, actual %d, %d, %s", tt.given, tt.expected.seasonNum, tt.expected.episodeNum, tt.expected.name, seasonNum, episodeNum, name)
			}
		})
	}
}

func Test_parseReleasedDate(t *testing.T) {
	var tests = []struct {
		given    string
		expected time.Time
	}{
		{"", time.Time{}},
		{"1", time.Time{}},
		{"Sun, Feb 3, 2005", time.Date(2005, time.February, 3, 0, 0, 0, 0, time.UTC)},
		{"Mon, Jan 2, 2006", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)},
		{"Tue, Mar 4, 2007", time.Date(2007, time.March, 4, 0, 0, 0, 0, time.UTC)},
		{"Wed, Apr 5, 2008", time.Date(2008, time.April, 5, 0, 0, 0, 0, time.UTC)},
		{"Thu, May 6, 2009", time.Date(2009, time.May, 6, 0, 0, 0, 0, time.UTC)},
		{"Fri, Jun 7, 2010", time.Date(2010, time.June, 7, 0, 0, 0, 0, time.UTC)},
		{"Sat, Jul 8, 2011", time.Date(2011, time.July, 8, 0, 0, 0, 0, time.UTC)},
		{"Sun, Aug 9, 2012", time.Date(2012, time.August, 9, 0, 0, 0, 0, time.UTC)},
		{"Mon, Sep 10, 2013", time.Date(2013, time.September, 10, 0, 0, 0, 0, time.UTC)},
		{"Tue, Oct 11, 2014", time.Date(2014, time.October, 11, 0, 0, 0, 0, time.UTC)},
		{"Wed, Nov 12, 2015", time.Date(2015, time.November, 12, 0, 0, 0, 0, time.UTC)},
		{"Thu, Dec 13, 2016", time.Date(2016, time.December, 13, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.given, func(t *testing.T) {
			actual := parseReleasedDate(tt.given)
			if actual != tt.expected {
				t.Errorf("parseReleasedDate(%s): expected %v, actual %v", tt.given, tt.expected, actual)
			}
		})
	}
}

func Test_parseStarRating(t *testing.T) {
	var tests = []struct {
		given    string
		expected float32
	}{
		{"1", 1},
		{"1.5", 1.5},
		{"10", 10},
		{"10.5", 0},
		{"0", 0},
		{"999", 0},
		{"hello", 0},
	}

	for _, tt := range tests {
		t.Run(tt.given, func(t *testing.T) {
			actual := parseStarRating(tt.given)
			if actual != tt.expected {
				t.Errorf("parseStarRating(%s): expected %f, actual %f", tt.given, tt.expected, actual)
			}
		})
	}
}

func Test_parseVoteCount(t *testing.T) {
	var tests = []struct {
		given    string
		expected int
	}{
		{" (148K)", 148000},
		{" (8K)", 8000},
		{" (12K)", 12000},
		{" (1)", 1},
		{" (10)", 10},
		{" (100)", 100},
		{" (1K)", 1000},
		{" (1.9K)", 1900},
	}

	for _, tt := range tests {
		t.Run(tt.given, func(t *testing.T) {
			actual := parseVoteCount(tt.given)
			if actual != tt.expected {
				t.Errorf("parseVoteCount(%s): expected %d, actual %d", tt.given, tt.expected, actual)
			}
		})
	}
}

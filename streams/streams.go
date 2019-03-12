package streams

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/pboehm/series/config"
	"log"
	"sort"
)

func BuildAbsoluteUrl(path string) string {
	return fmt.Sprintf("https://s.to%s", path)
}

type Series struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

func (s Series) AbsoluteLink() string {
	return BuildAbsoluteUrl(fmt.Sprintf("/serie/stream/%s", s.Link))
}

type Link struct {
	ID          int    `json:"id"`
	Link        string `json:"link"`
	Hoster      string `json:"hoster"`
	HosterTitle string `json:"hosterTitle"`
	Language    int    `json:"language"`
}

type Episode struct {
	ID          int     `json:"id"`
	Series      int     `json:"series"`
	Season      int     `json:"season"`
	Episode     int     `json:"episode"`
	German      string  `json:"german"`
	English     string  `json:"english"`
	Description string  `json:"description"`
	Links       []*Link `json:"links"`
	Language    struct {
		German  bool `json:"german"`
		English bool `json:"english"`
		GerSub  bool `json:"ger-sub"`
	} `json:"language"`
}

func (e Episode) LinksInLanguage(language int) []*Link {
	var links []*Link

	for _, link := range e.Links {
		if link.Language == language {
			links = append(links, link)
		}
	}

	return links
}

type AvailableSeriesResponse struct {
	Series []*Series `json:"series"`
}

type SeriesWithSeasonsResponse struct {
	Series  []*Series `json:"series"`
	Seasons []int     `json:"seasons"`
}

type SeriesWithEpisodesResponse struct {
	Series   []*Series  `json:"series"`
	Episodes []*Episode `json:"episodes"`
}

type Streams struct {
	Config config.Config
}

func (s *Streams) AvailableSeries() []*Series {
	header := req.Header{
		"Accept": "application/json",
	}
	param := req.QueryParam{
		"key":      s.Config.StreamsAPIToken,
		"extended": "0",
		"category": "0",
	}

	r, err := req.Get("https://s.to/api/v1/series/list", header, param)
	if err != nil {
		log.Fatal(err)
	}

	var parsedResponse AvailableSeriesResponse
	r.ToJSON(&parsedResponse)

	return parsedResponse.Series
}

func (s *Streams) Seasons(series *Series) []int {
	header := req.Header{
		"Accept": "application/json",
	}
	param := req.QueryParam{
		"key":    s.Config.StreamsAPIToken,
		"series": series.Id,
	}

	r, err := req.Get("https://s.to/api/v1/series/get", header, param)
	if err != nil {
		log.Fatal(err)
	}

	var parsedResponse SeriesWithSeasonsResponse
	r.ToJSON(&parsedResponse)

	seasons := parsedResponse.Seasons
	sort.Sort(sort.Reverse(sort.IntSlice(seasons)))
	return seasons
}

func (s *Streams) Episodes(series *Series, season int) []*Episode {
	header := req.Header{
		"Accept": "application/json",
	}
	param := req.QueryParam{
		"key":    s.Config.StreamsAPIToken,
		"series": series.Id,
		"season": season,
	}

	r, err := req.Get("https://s.to/api/v1/series/get", header, param)
	if err != nil {
		log.Fatal(err)
	}

	var parsedResponse SeriesWithEpisodesResponse
	r.ToJSON(&parsedResponse)

	return parsedResponse.Episodes
}

func (s *Streams) LinkUrl(link *Link) string {
	return BuildAbsoluteUrl(fmt.Sprintf("%s?key=%s", link.Link, s.Config.StreamsAPIToken))
}

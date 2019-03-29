package streams

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"github.com/pboehm/series/config"
	"log"
	"net/http"
	"sort"
)

type Series struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

func (s Series) AbsoluteLink() string {
	return absoluteUrl(fmt.Sprintf("/serie/stream/%s", s.Link))
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

type LoginResponse struct {
	Avatar   string `json:"avatar"`
	Session  string `json:"session"`
	Success  bool   `json:"success"`
	UID      int    `json:"uid"`
	UserLink string `json:"userlink"`
	Username string `json:"username"`
}

type Streams struct {
	Config   config.Config
	requests *req.Req
}

func NewStreams(config config.Config) *Streams {
	requests := req.New()
	requests.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &Streams{
		Config:   config,
		requests: requests,
	}
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

	r, err := s.requests.Get(absoluteUrl("/api/v1/series/list"), header, param)
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

	r, err := s.requests.Get(absoluteUrl("/api/v1/series/get"), header, param)
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

	r, err := s.requests.Get(absoluteUrl("/api/v1/series/get"), header, param)
	if err != nil {
		log.Fatal(err)
	}

	var parsedResponse SeriesWithEpisodesResponse
	r.ToJSON(&parsedResponse)

	return parsedResponse.Episodes
}

func (s *Streams) Login(email string, password string) (string, error) {
	body := req.Param{"email": email, "password": password}

	url := absoluteUrl(fmt.Sprintf("/api/v1/account/login?key=%s", s.Config.StreamsAPIToken))
	r, err := s.requests.Post(url, body)
	if err != nil {
		return "", err
	}

	var response LoginResponse
	if err = r.ToJSON(&response); err != nil {
		return "", err
	}

	return response.Session, nil
}

func (s *Streams) ResolveLink(linkId int, session string) (string, error) {
	header := req.Header{
		"Accept": "*/*",
		"Cookie": fmt.Sprintf("SSTOSESSION=%s", session),
	}

	r, err := s.requests.Head(s.LinkUrl(linkId), header)
	if err != nil {
		return "", err
	}

	response := r.Response()
	location := response.Header.Get("Location")
	if response.StatusCode != 301 || location == "" {
		return "", errors.New(fmt.Sprintf("could not resolve link: status=%d location=%s", response.StatusCode, location))
	}

	return location, nil
}

func (s *Streams) LinkUrl(linkId int) string {
	return absoluteUrl(fmt.Sprintf("/api/v1/stream/%d?key=%s", linkId, s.Config.StreamsAPIToken))
}

func absoluteUrl(path string) string {
	rev := func(s string) string {
		r := []rune(s)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	}

	return fmt.Sprintf("%s%s", rev("ot.s//:sptth"), path)
}

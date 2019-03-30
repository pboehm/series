package streams

import (
	"github.com/gin-gonic/gin"
	"github.com/pboehm/series/config"
	"sort"
)

type groupedSeriesResponse struct {
	Series   string          `json:"series"`
	Episodes []*LinkSetEntry `json:"episodes"`
}

type definedActionResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type API struct {
	Config              config.Config
	Jobs                *JobPool
	HtmlContent         func() []byte
	LinkSet             func() *LinkSet
	LinkSetRefresh      func()
	MarkWatched         func([]string) ([]string, []string)
	ExecuteLinkAction   func(config.StreamAction, *Identifier, int) *Job
	ExecuteGlobalAction func(config.StreamAction) *Job
}

//noinspection ALL
func (a *API) Run(listen string) error {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", a.HtmlContent())
	})
	r.GET("/api/links", func(c *gin.Context) {
		entriesReady := false
		//noinspection GoPreferNilSlice
		entries := []*LinkSetEntry{}

		linkSet := a.LinkSet()
		if linkSet != nil {
			entriesReady = true
			entries = linkSet.Entries()
		}

		c.JSON(200, gin.H{
			"ready": entriesReady,
			"links": entries,
		})
	})
	r.GET("/api/links/grouped", func(c *gin.Context) {
		entriesReady := false
		entries := map[string][]*LinkSetEntry{}

		linkSet := a.LinkSet()
		if linkSet != nil {
			entriesReady = true
			entries = linkSet.GroupedEntries()
		}

		//noinspection GoPreferNilSlice
		grouped := []groupedSeriesResponse{}
		for seriesName, entries := range entries {
			grouped = append(grouped, groupedSeriesResponse{
				Series:   seriesName,
				Episodes: entries,
			})
		}

		// sort grouped series so that series with newer episodes appear first
		sort.Slice(grouped, func(i, j int) bool {
			iEpisodes := grouped[i].Episodes
			jEpisodes := grouped[j].Episodes
			return iEpisodes[len(iEpisodes)-1].EpisodeId > jEpisodes[len(jEpisodes)-1].EpisodeId
		})

		c.JSON(200, gin.H{
			"ready": entriesReady,
			"links": grouped,
		})
	})
	r.POST("/api/links/refresh", func(c *gin.Context) {
		a.LinkSetRefresh()

		c.JSON(200, gin.H{
			"success": true,
		})
	})
	r.POST("/api/links/watched", func(c *gin.Context) {
		var episodeIds []string
		var err error

		if err = c.BindJSON(&episodeIds); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		successes, failures := a.MarkWatched(episodeIds)

		c.JSON(200, gin.H{
			"successes": successes,
			"failures":  failures,
		})
	})
	r.GET("/api/actions/global", func(c *gin.Context) {
		definedActions := []definedActionResponse{}
		for _, action := range a.Config.StreamsGlobalActions {
			definedActions = append(definedActions, definedActionResponse{
				Id:    action.Id,
				Title: action.Title,
			})
		}

		c.JSON(200, definedActions)
	})
	r.POST("/api/actions/global/:action", func(c *gin.Context) {
		action := config.StreamAction{}
		for _, a := range a.Config.StreamsGlobalActions {
			if a.Id == c.Param("action") {
				action = a
			}
		}

		if action.Id == "" {
			c.JSON(400, gin.H{
				"error": "action not found",
			})
			return
		}

		job := a.ExecuteGlobalAction(action)
		a.Jobs.Set(job.Id, job)
		go job.Run()

		c.JSON(200, job.Response())
	})
	r.GET("/api/actions/link", func(c *gin.Context) {
		definedActions := []definedActionResponse{}
		for _, action := range a.Config.StreamsLinkActions {
			definedActions = append(definedActions, definedActionResponse{
				Id:    action.Id,
				Title: action.Title,
			})
		}

		c.JSON(200, definedActions)
	})
	r.POST("/api/actions/link/:action/:linkId", func(c *gin.Context) {
		action := config.StreamAction{}
		for _, a := range a.Config.StreamsLinkActions {
			if a.Id == c.Param("action") {
				action = a
			}
		}

		if action.Id == "" {
			c.JSON(400, gin.H{
				"error": "action not found",
			})
			return
		}

		identifier, linkId, err := LinkIdentifierFromString(c.Param("linkId"))
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		job := a.ExecuteLinkAction(action, identifier, linkId)
		a.Jobs.Set(job.Id, job)
		go job.Run()

		c.JSON(200, job.Response())
	})
	r.GET("/api/actions/job/:job", func(c *gin.Context) {
		job, ok := a.Jobs.Get(c.Param("job"))
		if ok {
			c.JSON(200, job.Response())
		} else {
			c.JSON(400, gin.H{
				"error": "job not found",
			})
		}
	})

	return r.Run(listen)
}

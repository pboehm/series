package streams

import (
	"github.com/gin-gonic/gin"
)

type API struct {
	HtmlContent    func() []byte
	LinkSet        func() *LinkSet
	LinkSetRefresh func()
	MarkWatched    func([]string) ([]string, []string)
}

func (a *API) Run(listen string) error {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", a.HtmlContent())
	})
	r.GET("/api/links", func(c *gin.Context) {
		entriesReady := false
		var entries []*LinkSetEntry

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

		c.JSON(200, gin.H{
			"ready": entriesReady,
			"links": entries,
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

	return r.Run(listen)
}

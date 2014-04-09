package index

import (
    . "launchpad.net/gocheck"
    "github.com/pboehm/series/util"
    "github.com/pboehm/series/renamer"
    "testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
var _ = Suite(&MySuite{})

type MySuite struct{
    dir string
    index *SeriesIndex
    err error
}

func (s *MySuite) SetUpTest(c *C) {
    s.dir = c.MkDir()
    s.index, s.err = ParseSeriesIndex("data/seriesindex_example.xml")
}

func (s *MySuite) TestEnvironment(c *C) {
    c.Assert(util.PathExists(s.dir), Equals, true)
}

func (s *MySuite) TestIndexParsing(c *C) {
    c.Assert(s.index, NotNil)
    c.Assert(s.err, IsNil)
    c.Assert(s.index.SeriesList, HasLen, 4)
}

func (s *MySuite) TestSeriesLookupCache(c *C) {
    c.Assert(s.index.SeriesMap, HasLen, 6)
    c.Assert(s.index.SeriesMap["Community"], Equals, s.index.SeriesMap["Comm"])
}

func (s *MySuite) TestEpisodeLookupCache(c *C) {
    series := s.index.SeriesMap["Shameless US"]

    c.Assert(series.EpisodeMap, HasLen, 2)
    c.Assert(series.EpisodeMap["de"], HasLen, 8)
    c.Assert(series.EpisodeMap["en"], HasLen, 23)

    c.Assert(series.EpisodeMap["de"]["1_1"], Equals, "S01E01 - Pilot.avi")

    epi, exist := series.EpisodeMap["de"]["1_9"]
    c.Assert(exist, Equals, false)
    c.Assert(epi, Equals, "")
}

func (s *MySuite) TestEpisodeExistanceCheckWithExactSeriesName(c *C) {
    episode := renamer.Episode{ Series: "Shameless US", Season: 1, Episode: 1,
                                Language: "de" }
    c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)

    episode.Language = "en"
    c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)

    episode.Language = "fr"
    c.Assert(s.index.IsEpisodeInIndex(episode), Equals, false)

    episode.Season = 100
    c.Assert(s.index.IsEpisodeInIndex(episode), Equals, false)
}

func (s *MySuite) TestSeriesNameExistanceCheck(c *C) {
    c.Assert(s.index.SeriesNameInIndex("Shameless US"), Equals, "Shameless US")
    c.Assert(s.index.SeriesNameInIndex("shameless US"), Equals, "Shameless US")

    c.Assert(s.index.SeriesNameInIndex("Community"), Equals, "Community")
    c.Assert(s.index.SeriesNameInIndex("Comm"), Equals, "Community")
    c.Assert(s.index.SeriesNameInIndex("Unity"), Equals, "Community")
    c.Assert(s.index.SeriesNameInIndex("tvp community"), Equals, "Community")

    c.Assert(s.index.SeriesNameInIndex("tHE bIG bANG tHEORY"),
                                       Equals, "The Big Bang Theory")
}

package index

import (
	"github.com/pboehm/series/renamer"
	"github.com/pboehm/series/util"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&MySuite{})

type MySuite struct {
	dir   string
	index *SeriesIndex
	err   error
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
	c.Assert(s.index.seriesMap, HasLen, 6)
	c.Assert(s.index.seriesMap["Community"], Equals, s.index.seriesMap["Comm"])
}

func (s *MySuite) TestEpisodeLookupCache(c *C) {
	series := s.index.seriesMap["Shameless US"]

	c.Assert(series.episodeMap, HasLen, 2)
	c.Assert(series.episodeMap["de"], HasLen, 8)
	c.Assert(series.episodeMap["en"], HasLen, 23)

	c.Assert(series.episodeMap["de"]["1_1"], Equals, "S01E01 - Pilot.avi")

	epi, exist := series.episodeMap["de"]["1_9"]
	c.Assert(exist, Equals, false)
	c.Assert(epi, Equals, "")
}

func (s *MySuite) TestEpisodeExistanceCheckWithExactSeriesName(c *C) {
	episode := renamer.Episode{Series: "Shameless US", Season: 1, Episode: 1,
		Language: "de"}
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

func (s *MySuite) TestAddingValidEpisodeToIndex(c *C) {
	episode := renamer.Episode{Series: "Shameless US", Season: 1, Episode: 9,
		Name: "Testepisode", Extension: ".mkv",
		Language: "de"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
}

func (s *MySuite) TestAddingAlreadyExistingEpisodeToIndex(c *C) {
	episode := renamer.Episode{Series: "Shameless US", Season: 1, Episode: 1,
		Name: "Testepisode", Extension: ".mkv",
		Language: "de"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, ErrorMatches, "Episode already exists in Index")
	c.Assert(added, Equals, false)
}

func (s *MySuite) TestAddingEpisodeWithoutLanguageToSeriesWithSingleLang(c *C) {
	episode := renamer.Episode{Series: "The Big Bang Theory", Season: 6,
		Episode: 5, Name: "Testepisode",
		Extension: ".mkv"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
	c.Assert(episode.Language, Equals, "de")
}

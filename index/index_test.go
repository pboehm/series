package index

import (
	"github.com/pboehm/series/renamer"
	"github.com/pboehm/series/util"
	. "launchpad.net/gocheck"
	"path"
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

	c.Assert(series.languageMap, HasLen, 2)
	c.Assert(series.languageMap["de"].episodeMap, HasLen, 8)
	c.Assert(series.languageMap["en"].episodeMap, HasLen, 23)

	c.Assert(series.languageMap["de"].episodeMap["1_1"],
		Equals, "S01E01 - Pilot.avi")

	epi, exist := series.languageMap["de"].episodeMap["1_9"]
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

func (s *MySuite) TestEpisodeExistanceWithAllBefore(c *C) {
	episode := renamer.Episode{Series: "The Big Bang Theory", Season: 1,
		Episode: 1, Language: "de"}
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)

	episode = renamer.Episode{Series: "The Big Bang Theory", Season: 6,
		Episode: 0, Language: "de"}
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
}

func (s *MySuite) TestAddValidEpisodeToIndex(c *C) {
	episode := renamer.Episode{Series: "Shameless US", Season: 1, Episode: 9,
		Name: "Testepisode", Extension: ".mkv", Language: "de"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
}

func (s *MySuite) TestAddAlreadyExistingEpisodeToIndex(c *C) {
	episode := renamer.Episode{Series: "Shameless US", Season: 1, Episode: 1,
		Name: "Testepisode", Extension: ".mkv", Language: "de"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, ErrorMatches, "episode already exists in index")
	c.Assert(added, Equals, false)
}

func (s *MySuite) TestAddEpisodeWithoutLanguageToSeriesWithSingleLang(c *C) {
	episode := renamer.Episode{Series: "The Big Bang Theory", Season: 6,
		Episode: 5, Name: "Testepisode", Extension: ".mkv"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
	c.Assert(episode.Language, Equals, "de")
}

func (s *MySuite) TestAddEpisodeWithoutLanguageToSeriesWithMultiLang(c *C) {
	// this episode has already been watched in "en" so it has to be "de"
	episode := renamer.Episode{Series: "Shameless US", Season: 1,
		Episode: 9, Name: "Testepisode", Extension: ".mkv"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
	c.Assert(episode.Language, Equals, "de")
}

func (s *MySuite) TestAddEpisodeWithoutLanguageToSeriesWithMultiLangPrev(c *C) {
	// this episode hasn't been watched in any language but there is one series
	// where the previous episode exists
	episode := renamer.Episode{Series: "Shameless US", Season: 2,
		Episode: 12, Name: "Testepisode", Extension: ".mkv"}

	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, false)
	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)
	c.Assert(episode.Language, Equals, "en")
}

func (s *MySuite) TestAddEpisodeWithCrappyFileInfos(c *C) {
	episode := renamer.Episode{
		Series: "tvs tbbt dd51 ded dl 18p ithd avc", Season: 6,
		Episode: 5, Name: "", Extension: ".mkv"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(added, Equals, false)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "series does not exist in index")
}

type mockExtractor struct {
	names []string
}

func (m mockExtractor) Names(*renamer.Episode) ([]string, error) {
	return m.names, nil
}

func (s *MySuite) TestSeriesNameExtractor(c *C) {
	s.index.AddExtractor(mockExtractor{names: []string{
		"Should-Also-Not-Exist", "Shameless US"}})

	episode := renamer.Episode{Series: "this-should-not-exist",
		Season: 1, Episode: 9, Name: "Testepisode", Extension: ".mkv",
		Language: "de"}

	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(episode.Series, Equals, "Shameless US")
}

func (s *MySuite) TestWriteIndexToFile(c *C) {
	episode := renamer.Episode{Series: "Shameless US", Season: 1, Episode: 9,
		Name: "Testepisode", Extension: ".mkv", Language: "de"}

	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, false)

	// add episode
	added, err := s.index.AddEpisode(&episode)
	c.Assert(err, IsNil)
	c.Assert(added, Equals, true)
	c.Assert(s.index.IsEpisodeInIndex(episode), Equals, true)

	// dump it
	dest := path.Join(s.dir, "seriesindex_dump.xml")
	s.index.WriteToFile(dest)
	c.Assert(util.PathExists(dest), Equals, true)

	// parse it back and make assertions
	index, err := ParseSeriesIndex(dest)
	c.Assert(err, IsNil)
	c.Assert(index.IsEpisodeInIndex(episode), Equals, true)
}

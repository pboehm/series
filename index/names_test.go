package index

import (
	"github.com/pboehm/series/renamer"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"testing"
)

func TestExtractor(t *testing.T) { TestingT(t) }

var _ = Suite(&ExtractorSuite{})

type ExtractorSuite struct {
	dir      string
	fixtures map[string]EpisodeFixture
}

type EpisodeFixture struct {
	path  string
	dir   bool
	files map[string]string
}

func createFile(path string, content string) {
	_ = ioutil.WriteFile(path, []byte(content), 0644)
}

func (s *ExtractorSuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	s.fixtures = map[string]EpisodeFixture{
		"crmi": {
			"Criminal.Minds.S01E01.Testtest.mkv",
			false, map[string]string{}},
		"rules_of_engagement": {
			"RoEG8p.713/Rules.of.Engagement.S07E13.100th.GERMAN.DL.DUBBED/",
			true, map[string]string{
				"tvp-egagement-s07e13-1080p.mkv": "abcksfvfddvhfjv",
				"tvp-egagement-s07e13-1080p.nfo": "abc",
			}},
	}

	for key, fixture := range s.fixtures {
		if fixture.dir {
			os.MkdirAll(s.FileWithPath(key), 0700)
			for file, content := range fixture.files {
				createFile(path.Join(s.FileWithPath(key), file), content)
			}

		} else {
			createFile(s.FileWithPath(key), "")
		}
	}
}

func (s *ExtractorSuite) FileWithPath(key string) string {
	return path.Join(s.dir, s.fixtures[key].path)
}

/////////////////////////
// Start of Test function

func (s *ExtractorSuite) TestEpisodePossibleSeriesNamesFromDirectory(c *C) {
	extractor := FilesystemExtractor{}
	episode, _ := renamer.CreateEpisodeFromPath(
		path.Dir(s.FileWithPath("rules_of_engagement")))

	names, err := extractor.Names(episode)
	c.Assert(names, DeepEquals, []string{
		"RoEG8p", "Rules of Engagement", "tvp egagement"})
	c.Assert(err, IsNil)
}

func (s *ExtractorSuite) TestEpisodePossibleSeriesNamesFromFile(c *C) {
	extractor := FilesystemExtractor{}
	episode, _ := renamer.CreateEpisodeFromPath(s.FileWithPath("crmi"))

	names, err := extractor.Names(episode)
	c.Assert(names, DeepEquals, []string{"Criminal Minds"})
	c.Assert(err, IsNil)
}

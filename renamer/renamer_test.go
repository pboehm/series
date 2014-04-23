package renamer

import (
	"github.com/pboehm/series/util"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&MySuite{})

type MySuite struct {
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

func (s *MySuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	s.fixtures = map[string]EpisodeFixture{
		"crmi": {
			"Criminal.Minds.S01E01.Testtest.mkv",
			false, map[string]string{}},
		"crmi_no_video": {
			"Criminal.Minds.S01E01.Testtest.pdf",
			false, map[string]string{}},
		"chuck1": {
			"Chuck.S01E01.Dies.ist.ein.Test.German.Dubbed.BLURAYRiP.mkv",
			false, map[string]string{}},
		"chuck2": {
			"chuck.512.hdtv-lol.avi",
			false, map[string]string{}},
		"chuck3": {
			"chuck.1212.hdtv-lol.avi",
			false, map[string]string{}},
		"chuck4": {
			"chuck.5x12.hdtv-lol.avi",
			false, map[string]string{}},
		"unknown_series": {
			"5x12.avi",
			false, map[string]string{}},
		"royal": {
			"Royal.Pains.S02E10.Beziehungsbeschwerden.GERMAN.DUBBED.avi",
			false, map[string]string{}},
		"ncis": {
			"NCIS.S11E13.Gueterzug.nach.Miami.GERMAN.DUBBED.DL.720p.WebHD.h264-euHD.mkv",
			false, map[string]string{}},
		"flpo1": {
			"Flashpoint.S04E04.Getruebte.Erinnerungen.German.Dubbed.avi",
			false, map[string]string{}},
		"flpo2": {
			"flpo.404.Die.German.Erinnerungen.German.Dubbed.BLURAYRiP.avi",
			false, map[string]string{}},
		"csi": {
			"sof-csi.ny.s07e20.avi",
			false, map[string]string{}},

		// sample illegal data
		"illegal1": {
			".DS_Store", false, map[string]string{}},
		"illegal2": {
			"Test", false, map[string]string{}},

		// sample directory data
		"crmi_dir": {
			"Criminal.Minds.S01E01.Testtest",
			true, map[string]string{
				"episode.mkv": "abcksfvfddvhfjvdhfvjdhfv",
				"sample.mkv":  "probablyshorter",
				"episode.sub": "ttttttttttttttttttttttttttttttttttttttttttt",
			}},
		"himym": {
			"HMM8p.909",
			true, map[string]string{
				"How.I.Met.Your.Mother.S09E09.Platonish.1080p.WEB-DL.DD5.mkv": "abcksfvfddvhfjvdhfvjdhfv",
			}},
		"himym_not_matching": {
			"HIMYM.909",
			true, map[string]string{
				"How.I.Met.Your.Mother.S09E10.Platonish.1080p.WEB-DL.DD5.mkv": "abcksfvfddvhfjvdhfvjdhfv",
			}},
		"rules_of_engagement": {
			"RoEG8p.713/Rules.of.Engagement.S07E13.100th.GERMAN.DL.DUBBED/",
			true, map[string]string{
				"tvp-egagement-s07e13-1080p.mkv": "abcksfvfddvhfjv",
				"tvp-egagement-s07e13-1080p.nfo": "abc",
			}},
		"chuck1_dir": {
			"Chuck.S01E01.Dies.ist.ein.Test.German.Dubbed.BLURAYRiP",
			true, map[string]string{}},
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

func (s *MySuite) FileWithPath(key string) string {
	return path.Join(s.dir, s.fixtures[key].path)
}

func (s *MySuite) TestEnvironment(c *C) {
	c.Assert(util.PathExists(s.dir), Equals, true)
	c.Assert(util.PathExists(s.FileWithPath("royal")), Equals, true)
	c.Assert(util.PathExists(s.FileWithPath("crmi_dir")), Equals, true)
	c.Assert(util.PathExists(
		path.Join(s.FileWithPath("crmi_dir"), "episode.mkv")), Equals, true)
}

func (s *MySuite) TestEpisodeInformationCleanup(c *C) {
	c.Assert(CleanEpisodeInformation("Criminal.Minds"),
		Equals, "Criminal Minds")
	c.Assert(CleanEpisodeInformation(".Criminal.Minds "),
		Equals, "Criminal Minds")
	c.Assert(CleanEpisodeInformation("GERMAN.DUBBED.DL.WEB-DL"),
		Equals, "GERMAN DUBBED DL WEB DL")
}

func (s *MySuite) TestInterestingFiles(c *C) {
	TestData := map[string]bool{
		"Criminal.Minds.S01E01.Testtest":                         true,
		"Chuck.S01E01.Dies.ist.ein.Test.German.Dubbed.BLURAYRiP": true,
		"chuck.512.hdtv-lol.avi":                                 true,
		"chuck.1212.hdtv-lol.avi":                                true,
		"chuck.5x12.hdtv-lol.avi":                                true,
		"5x12.avi":                                               true,
		"Royal.Pains.S02E10.Beziehungsbeschwerden.GERMAN.DUBBED.avi":     true,
		"Flashpoint.S04E04.Getruebte.Erinnerungen.German.Dubbed.avi":     true,
		"sof-csi.ny.s07e20.avi":                                          true,
		"flpo.404.Die.German.Erinnerungen.German.Dubbed.WEB-DL.XViD.avi": true,

		// sample illegal data
		".DS_Store": false,
		"Test":      false,
	}

	for key, val := range TestData {
		c.Assert(IsInterestingDirEntry(key), Equals, val,
			Commentf("IsInterestingDirEntry(%s) should be %v", key, val))
	}
}

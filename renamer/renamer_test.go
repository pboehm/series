package renamer

import (
    . "launchpad.net/gocheck"
    "github.com/pboehm/series/util"
    "testing"
    "path"
    "os"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
var _ = Suite(&MySuite{})

type MySuite struct{
    dir string
    files map[string]string
}

func (s *MySuite) SetUpTest(c *C) {
    s.dir = c.MkDir()
    s.files = map[string]string {
        "crmi":   "Criminal.Minds.S01E01.Testtest.mkv",
        "chuck1": "Chuck.S01E01.Dies.ist.ein.Test.German.Dubbed.BLURAYRiP.mkv",
        "chuck2": "chuck.512.hdtv-lol.avi",
        "chuck3": "chuck.1212.hdtv-lol.avi",
        "chuck4": "chuck.5x12.hdtv-lol.avi",
        "unknown_series": "5x12.avi",
        "royal": "Royal.Pains.S02E10.Beziehungsbeschwerden.GERMAN.DUBBED.avi",
        "flpo1": "Flashpoint.S04E04.Getruebte.Erinnerungen.German.Dubbed.avi",
        "flpo2": "flpo.404.Die.German.Erinnerungen.German.Dubbed.BLURAYRiP.avi",
        "csi": "sof-csi.ny.s07e20.avi",

        // sample illegal data
        "illegal1": ".DS_Store",
        "illegal2": "Test",
    }

    for key, _ := range s.files {
        file, _ := os.Create(s.FileWithPath(key))
        file.Close()
    }
}

func (s *MySuite) FileWithPath(key string) string {
    return path.Join(s.dir, s.files[key])
}

func (s *MySuite) TestEnvironment(c *C) {
    c.Assert(util.PathExists(s.dir), Equals, true)
    c.Assert(util.PathExists(s.FileWithPath("royal")), Equals, true)
}

func (s *MySuite) TestEpisodeInformationCleanup(c *C) {
    c.Assert(CleanEpisodeInformation("Criminal.Minds"),
        Equals, "Criminal Minds")
    c.Assert(CleanEpisodeInformation(".Criminal.Minds "),
        Equals, "Criminal Minds")
}

func (s *MySuite) TestInterestingFiles(c *C) {
    TestData := map[string]bool {
        "Criminal.Minds.S01E01.Testtest": true,
        "Chuck.S01E01.Dies.ist.ein.Test.German.Dubbed.BLURAYRiP": true,
        "chuck.512.hdtv-lol.avi": true,
        "chuck.1212.hdtv-lol.avi": true,
        "chuck.5x12.hdtv-lol.avi": true,
        "5x12.avi": true,
        "Royal.Pains.S02E10.Beziehungsbeschwerden.GERMAN.DUBBED.avi": true,
        "Flashpoint.S04E04.Getruebte.Erinnerungen.German.Dubbed.avi": true,
        "sof-csi.ny.s07e20.avi": true,
        "flpo.404.Die.German.Erinnerungen.German.Dubbed.WEB-DL.XViD.avi": true,

        // sample illegal data
        ".DS_Store": false,
        "Test": false,
    }

    for key, val := range TestData {
        c.Assert(IsInterestingDirEntry(key), Equals, val,
                Commentf("IsInterestingDirEntry(%s) should be %v", key, val))
    }
}

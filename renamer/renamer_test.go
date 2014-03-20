package renamer

import (
    . "launchpad.net/gocheck"
    "testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}
var _ = Suite(&MySuite{})

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
        c.Check(IsInterestingDirEntry(key), Equals, val,
                Commentf("IsInterestingDirEntry(%s) should be %v", key, val))
    }

}

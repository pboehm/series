package renamer

import (
    . "launchpad.net/gocheck"
)

func (s *MySuite) TestEpisodeStruct(c *C) {
    episode := Episode{season: 21, episode: 12, name: "Testepisode" }
    c.Check(episode.season, Equals, 21)
    c.Check(episode.CleanedFileName(), Equals, "S21E12 - Testepisode")

    episode = Episode{season: 1, episode: 1, name: "Testepisode" }
    c.Check(episode.CleanedFileName(), Equals, "S01E01 - Testepisode")
}

func (s *MySuite) TestEpisodeExtractionFromFile(c *C) {
    episode, err := CreateEpisodeFromPath(s.FileWithPath("crmi"))

    c.Check(err, IsNil)
    c.Check(episode, FitsTypeOf, new(Episode))
    c.Check(episode.season, Equals, 1)
    c.Check(episode.episode, Equals, 1)
    c.Check(episode.name, Equals, ".Testtest.mkv")
    c.Check(episode.series, Equals, "Criminal.Minds.")
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
        c.Check(IsInterestingDirEntry(key), Equals, val,
                Commentf("IsInterestingDirEntry(%s) should be %v", key, val))
    }
}

package renamer

import (
    . "launchpad.net/gocheck"
)

func (s *MySuite) TestEpisodeStruct(c *C) {
    episode := Episode{season: 21, episode: 12, name: "Testepisode" }
    c.Assert(episode.season, Equals, 21)
    c.Assert(episode.CleanedFileName(), Equals, "S21E12 - Testepisode")

    episode = Episode{season: 1, episode: 1, name: "Testepisode" }
    c.Assert(episode.CleanedFileName(), Equals, "S01E01 - Testepisode")
}

func (s *MySuite) TestEpisodeExtractionFromFile(c *C) {
    episode, err := CreateEpisodeFromPath(s.FileWithPath("crmi"))

    c.Assert(err, IsNil)
    c.Assert(episode, Not(IsNil))
    c.Assert(episode, FitsTypeOf, new(Episode))
    c.Assert(episode.season, Equals, 1)
    c.Assert(episode.episode, Equals, 1)
    c.Assert(episode.name, Equals, "Testtest")
    c.Assert(episode.series, Equals, "Criminal Minds")
    c.Assert(episode.CleanedFileName(), Equals, "S01E01 - Testtest.mkv")
}

func (s *MySuite) TestEpisodeExtractionFromInvalidFile(c *C) {
    _, err := CreateEpisodeFromPath(s.FileWithPath("illegal1"))

    c.Assert(err, ErrorMatches, "Supplied episode has no series information")
}

func (s *MySuite) TestEpisodeExtractionFromNotExistingFile(c *C) {
    _, err := CreateEpisodeFromPath("/should/not/exist")

    c.Assert(err, ErrorMatches, "Supplied episode does not exist")
}

func (s *MySuite) TestEpisodeTrashwordRemoval(c *C) {
    episode, _ := CreateEpisodeFromPath(s.FileWithPath("flpo2"))

    c.Assert(episode.name, Equals,
        "Die German Erinnerungen German Dubbed BLURAYRiP")
    episode.RemoveTrashwords()
    c.Assert(episode.name, Equals, "Die German Erinnerungen")
}

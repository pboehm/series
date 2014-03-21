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
    c.Check(episode.name, Equals, "Testtest mkv")
    c.Check(episode.series, Equals, "Criminal Minds")
}

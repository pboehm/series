package renamer

import (
    . "launchpad.net/gocheck"
    "path"
    "github.com/pboehm/series/util"
)

func (s *MySuite) TestEpisodeStruct(c *C) {
    episode := Episode{Season: 21, Episode: 12, Name: "Testepisode" }
    c.Assert(episode.Season, Equals, 21)
    c.Assert(episode.CleanedFileName(), Equals, "S21E12 - Testepisode")

    episode = Episode{Season: 1, Episode: 1, Name: "Testepisode" }
    c.Assert(episode.CleanedFileName(), Equals, "S01E01 - Testepisode")
}

func (s *MySuite) TestEpisodeExtractionFromFile(c *C) {
    episode, err := CreateEpisodeFromPath(s.FileWithPath("crmi"))

    c.Assert(err, IsNil)
    c.Assert(episode, Not(IsNil))
    c.Assert(episode, FitsTypeOf, new(Episode))
    c.Assert(episode.Season, Equals, 1)
    c.Assert(episode.Episode, Equals, 1)
    c.Assert(episode.Name, Equals, "Testtest")
    c.Assert(episode.Series, Equals, "Criminal Minds")
    c.Assert(episode.Episodefile, Equals, s.FileWithPath("crmi"))
    c.Assert(episode.CleanedFileName(), Equals, "S01E01 - Testtest.mkv")
}

func (s *MySuite) TestEpisodeExtractionFromDirectory(c *C) {
    episode, err := CreateEpisodeFromPath(s.FileWithPath("crmi_dir"))

    c.Assert(err, IsNil)
    c.Assert(episode, Not(IsNil))
    c.Assert(episode, FitsTypeOf, new(Episode))
    c.Assert(episode.Season, Equals, 1)
    c.Assert(episode.Episode, Equals, 1)
    c.Assert(episode.Name, Equals, "Testtest")
    c.Assert(episode.Series, Equals, "Criminal Minds")
    c.Assert(episode.Episodefile, Equals,
        path.Join(s.FileWithPath("crmi_dir"), "episode.mkv"))
    c.Assert(episode.CleanedFileName(), Equals, "S01E01 - Testtest.mkv")
    c.Assert(episode.CanBeRenamed(), Equals, true)
}

func (s *MySuite) TestEpisodeExtractionFromDirectoryWithBetterFile(c *C) {
    episode, err := CreateEpisodeFromPath(s.FileWithPath("himym"))

    c.Assert(err, IsNil)
    c.Assert(episode, Not(IsNil))
    c.Assert(episode, FitsTypeOf, new(Episode))
    c.Assert(episode.Season, Equals, 9)
    c.Assert(episode.Episode, Equals, 9)
    c.Assert(episode.Name, Equals, "Platonish")
    c.Assert(episode.Series, Equals, "How I Met Your Mother")
    c.Assert(episode.Episodefile, Equals,
        path.Join(s.FileWithPath("himym"),
            "How.I.Met.Your.Mother.S09E09.Platonish.1080p.WEB-DL.DD5.mkv"))
    c.Assert(episode.CleanedFileName(), Equals, "S09E09 - Platonish.mkv")
    c.Assert(episode.CanBeRenamed(), Equals, true)
}

func (s *MySuite) TestEpisodeExtractionFromDirectoryWithInvalidBetterFile(c *C) {
    episode, err := CreateEpisodeFromPath(s.FileWithPath("himym_not_matching"))

    c.Assert(err, IsNil)
    c.Assert(episode, Not(IsNil))
    c.Assert(episode, FitsTypeOf, new(Episode))
    c.Assert(episode.Season, Equals, 9)
    c.Assert(episode.Episode, Equals, 9)
    c.Assert(episode.Name, Equals, "")
    c.Assert(episode.Series, Equals, "HIMYM")
    c.Assert(episode.CanBeRenamed(), Equals, false)
}

func (s *MySuite) TestEpisodeThatShouldntBeRenamable(c *C) {
    episode, _ := CreateEpisodeFromPath(s.FileWithPath("unknown_series"))
    c.Assert(episode.CanBeRenamed(), Equals, false)
    c.Assert(episode.HasValidEpisodeName(), Equals, false)
}

func (s *MySuite) TestEpisodeExtractionFromInvalidFile(c *C) {
    _, err := CreateEpisodeFromPath(s.FileWithPath("illegal1"))
    c.Assert(err, ErrorMatches, "Supplied episode has no series information")
}

func (s *MySuite) TestEpisodeExtractionFromNoVideoFile(c *C) {
    _, err := CreateEpisodeFromPath(s.FileWithPath("crmi_no_video"))
    c.Assert(err, ErrorMatches, "No videofile available")
}

func (s *MySuite) TestEpisodeExtractionFromDirectoryWithoutVideoFile(c *C) {
    _, err := CreateEpisodeFromPath(s.FileWithPath("chuck1_dir"))
    c.Assert(err, ErrorMatches, "No videofile available")
}

func (s *MySuite) TestEpisodeExtractionFromNotExistingFile(c *C) {
    _, err := CreateEpisodeFromPath("/should/not/exist")
    c.Assert(err, ErrorMatches, "Supplied episode does not exist")
}

func (s *MySuite) TestEpisodeTrashwordRemoval(c *C) {
    episode, _ := CreateEpisodeFromPath(s.FileWithPath("flpo2"))

    c.Assert(episode.Name, Equals,
        "Die German Erinnerungen German Dubbed BLURAYRiP")
    episode.RemoveTrashwords()
    c.Assert(episode.Name, Equals, "Die German Erinnerungen")
}

func (s *MySuite) TestEpisodeTrashwordRemovalSkipAfterTwoPurges(c *C) {
    episode, _ := CreateEpisodeFromPath(s.FileWithPath("ncis"))
    episode.RemoveTrashwords()

    c.Assert(episode.Name, Equals, "Gueterzug nach Miami")
}

func (s *MySuite) TestEpisodeRenamingEpisodeFile(c *C) {
    episode, _ := CreateEpisodeFromPath(s.FileWithPath("chuck1"))
    c.Assert(episode.CanBeRenamed(), Equals, true)
    episode.Rename(s.dir)
    c.Assert(util.PathExists(path.Join(s.dir, episode.CleanedFileName())),
             Equals, true)
    c.Assert(util.PathExists(s.FileWithPath("chuck1")), Equals, false)
}

func (s *MySuite) TestEpisodeRenamingFromDir(c *C) {
    episode, _ := CreateEpisodeFromPath(s.FileWithPath("crmi_dir"))
    c.Assert(episode.CanBeRenamed(), Equals, true)
    episode.Rename(s.dir)
    c.Assert(util.PathExists(path.Join(s.dir, episode.CleanedFileName())),
             Equals, true)
    c.Assert(util.PathExists(s.FileWithPath("crmi_dir")), Equals, false)
}


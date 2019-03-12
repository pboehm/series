package renamer

import (
	"errors"
	"fmt"
	"github.com/pboehm/series/util"
	"os"
	GlobalPath "path"
	"path/filepath"
	"regexp"
	"strconv"
)

func CreateEpisodeFromPath(path string) (*Episode, error) {
	episode := new(Episode)

	if !util.PathExists(path) {
		return episode, errors.New("supplied episode does not exist")
	}

	basename := filepath.Base(path)
	if !IsInterestingDirEntry(basename) {
		return episode, errors.New("supplied episode has no series information")
	}

	episode.Path = path
	episode.EpisodeFile = path
	if util.IsDirectory(path) {
		episodeFile, err := FindBiggestVideoFile(path)

		if err != nil {
			return episode, err
		}
		episode.EpisodeFile = episodeFile
	}

	if !HasVideoFileEnding(episode.EpisodeFile) {
		return episode, errors.New("no video file available")
	}

	information := ExtractEpisodeInformation(basename)
	episode.Season, _ = strconv.Atoi(information["season"])
	episode.Episode, _ = strconv.Atoi(information["episode"])

	episode.Series = CleanEpisodeInformation(information["series"])
	episode.Extension = GlobalPath.Ext(episode.EpisodeFile)

	name := information["episodename"]
	if util.IsFile(path) {
		name = name[:len(name)-len(episode.Extension)]
	}
	episode.Name = CleanEpisodeInformation(name)

	episode.ExtractLanguage()

	return episode, nil
}

type Episode struct {
	Season, Episode                                      int
	Name, Series, Extension, EpisodeFile, Path, Language string
}

func (e *Episode) CleanedFileName() string {
	return fmt.Sprintf("S%02dE%02d - %s%s",
		e.Season, e.Episode, e.Name, e.Extension)
}

func (e *Episode) HasValidEpisodeName() bool {
	return e.Name != ""
}

func (e *Episode) SetDefaultEpisodeName() {
	e.Name = fmt.Sprintf("Episode %02d", e.Episode)
}

func (e *Episode) CanBeRenamed() bool {
	return e.HasValidEpisodeName() && util.IsFile(e.EpisodeFile)
}

func (e *Episode) ExtractLanguage() {
	pattern := regexp.MustCompile("(?i)German")
	if pattern.Match([]byte(e.Name)) {
		e.Language = "de"
	}
}

func (e *Episode) RemoveTrashWords() {
	e.Name = ApplyTrashWordsOnString(e.Name)
}

func (e *Episode) Rename(destPath string) error {
	if !e.CanBeRenamed() {
		return errors.New(
			"this episode couldn't be renamed as it has some problems")
	}

	needCleanup := false
	if util.IsDirectory(e.Path) {
		needCleanup = true
	}

	dest := GlobalPath.Join(destPath, e.CleanedFileName())

	err := os.Rename(e.EpisodeFile, dest)
	if err != nil {
		return err
	}

	if needCleanup {
		return os.RemoveAll(e.Path)
	}

	return nil
}

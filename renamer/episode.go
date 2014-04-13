package renamer

import (
	"errors"
	"fmt"
	"github.com/pboehm/series/util"
	"os"
	GlobalPath "path"
	"path/filepath"
	"strconv"
	"regexp"
)

func CreateEpisodeFromPath(path string) (*Episode, error) {
	episode := new(Episode)

	if !util.PathExists(path) {
		return episode, errors.New("Supplied episode does not exist")
	}

	basename := filepath.Base(path)
	if !IsInterestingDirEntry(basename) {
		return episode, errors.New("Supplied episode has no series information")
	}

	episode.Path = path
	episode.Episodefile = path
	if util.IsDirectory(path) {
		episodefile, err := FindBiggestVideoFile(path)

		if err != nil {
			return episode, err
		}
		episode.Episodefile = episodefile
	}

	if !HasVideoFileEnding(episode.Episodefile) {
		return episode, errors.New("No videofile available")
	}

	information := ExtractEpisodeInformation(basename)
	episode.Season, _ = strconv.Atoi(information["season"])
	episode.Episode, _ = strconv.Atoi(information["episode"])

	episode.Series = CleanEpisodeInformation(information["series"])
	episode.Extension = GlobalPath.Ext(episode.Episodefile)

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
	Name, Series, Extension, Episodefile, Path, Language string
}

func (self *Episode) CleanedFileName() string {
	return fmt.Sprintf("S%02dE%02d - %s%s",
		self.Season, self.Episode, self.Name, self.Extension)
}

func (self *Episode) HasValidEpisodeName() bool {
	return self.Name != ""
}

func (self *Episode) SetDefaultEpisodeName() {
	self.Name = fmt.Sprintf("Episode %02d", self.Episode)
}

func (self *Episode) CanBeRenamed() bool {
	return self.HasValidEpisodeName() && util.IsFile(self.Episodefile)
}

func (self *Episode) ExtractLanguage() {
    pattern := regexp.MustCompile("(?i)German")
    if pattern.Match([]byte(self.Name)) {
        self.Language = "de"
    }
}

func (self *Episode) RemoveTrashwords() {
	self.Name = ApplyTrashwordsOnString(self.Name)
}

func (self *Episode) FindBetterInformation() {
	// check if the Episodefile contains more informations than we already have
	// and set these informations
	if util.IsDirectory(self.Path) {
		subepisode, suberr := CreateEpisodeFromPath(self.Episodefile)
		if suberr == nil {
			subepisode.RemoveTrashwords()

			if self.Season == subepisode.Season &&
				self.Episode == subepisode.Episode {

				if len(subepisode.Series) > len(self.Series) {
					self.Series = subepisode.Series
				}

				if len(subepisode.Name) > len(self.Name) {
					self.Name = subepisode.Name
				}
			}
		}
	}
}

func (self *Episode) Rename(dest_path string) error {
	if !self.CanBeRenamed() {
		return errors.New(
			"This episode couldn't be renamed as it has some problems")
	}

	need_cleanup := false
	if util.IsDirectory(self.Path) {
		need_cleanup = true
	}

	dest := GlobalPath.Join(dest_path, self.CleanedFileName())

	err := os.Rename(self.Episodefile, dest)
	if err != nil {
		return err
	}

	if need_cleanup {
		return os.RemoveAll(self.Path)
	}

	return nil
}

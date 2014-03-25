package renamer

import (
    "fmt"
    "github.com/pboehm/series/util"
    "errors"
    "path/filepath"
    GlobalPath "path"
    "strconv"
    "os"
)

func CreateEpisodeFromPath(path string) (*Episode, error) {
    episode := new(Episode)

    if ! util.PathExists(path) {
        return episode, errors.New("Supplied episode does not exist")
    }

    basename := filepath.Base(path)
    if ! IsInterestingDirEntry(basename) {
        return episode, errors.New("Supplied episode has no series information")
    }

    episode.path = path
    episode.episodefile = path
    if util.IsDirectory(path) {
        episodefile, err := FindBiggestVideoFile(path)

        if err != nil {
            return episode, err
        }
        episode.episodefile = episodefile
    }

    if ! HasVideoFileEnding(episode.episodefile) {
        return episode, errors.New("No videofile available")
    }

    information := ExtractEpisodeInformation(basename)
    episode.season, _  = strconv.Atoi(information["season"])
    episode.episode, _ = strconv.Atoi(information["episode"])

    episode.series = CleanEpisodeInformation(information["series"])
    episode.extension = GlobalPath.Ext(episode.episodefile)

    name := information["episodename"]
    if util.IsFile(path) {
        name = name[:len(name) - len(episode.extension)]
    }
    episode.name = CleanEpisodeInformation(name)

    return episode, nil
}

type Episode struct {
    season, episode int
    name, series, extension, episodefile, path string
}

func (self *Episode) CleanedFileName() string {
    return fmt.Sprintf("S%02dE%02d - %s%s",
                self.season, self.episode, self.name, self.extension)
}

func (self *Episode) HasValidEpisodeName() bool {
    return self.name != ""
}

func (self *Episode) SetDefaultEpisodeName() {
    self.name = fmt.Sprintf("Episode %02d", self.episode)
}

func (self *Episode) CanBeRenamed() bool {
    return self.HasValidEpisodeName() && util.IsFile(self.episodefile)
}

func (self *Episode) RemoveTrashwords() {
    self.name = ApplyTrashwordsOnString(self.name)
}

func (self *Episode) Rename(dest_path string) error {
    if ! self.CanBeRenamed() {
        return errors.New(
            "This episode couldn't be renamed as it has some problems")
    }

    need_cleanup := false
    if util.IsDirectory(self.path) {
        need_cleanup = true
    }

    dest := GlobalPath.Join(dest_path, self.CleanedFileName())

    err := os.Rename(self.episodefile, dest)
    if err != nil { return err }

    if need_cleanup {
        return os.RemoveAll(self.path)
    }

    return nil
}

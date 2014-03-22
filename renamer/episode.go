package renamer

import (
    "fmt"
    "github.com/pboehm/series/util"
    "errors"
    "path/filepath"
    GlobalPath "path"
    "strconv"
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

    information := ExtractEpisodeInformation(basename)
    episode.season, _  = strconv.Atoi(information["season"])
    episode.episode, _ = strconv.Atoi(information["episode"])

    episode.series = CleanEpisodeInformation(information["series"])

    name := information["episodename"]
    if util.IsFile(path) {
        episode.extension = GlobalPath.Ext(basename)
        name = name[:len(name) - len(episode.extension)]
    }

    episode.name   = CleanEpisodeInformation(name)

    return episode, nil
}

type Episode struct {
    season, episode int
    name, series, extension string
}

func (self Episode) CleanedFileName() string {
    return fmt.Sprintf("S%02dE%02d - %s%s",
                self.season, self.episode, self.name, self.extension)
}

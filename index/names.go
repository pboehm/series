package index

import (
	"fmt"
	"github.com/pboehm/series/renamer"
	"github.com/pboehm/series/util"
	"os/exec"
	"path"
	"strings"
)

type SeriesNameExtractor interface {
	Names(*renamer.Episode) ([]string, error)
}

type FilesystemExtractor struct{}

func (f FilesystemExtractor) Names(episode *renamer.Episode) ([]string, error) {
	possibilities := []string{episode.Series}

	if util.IsDirectory(episode.Path) {
		// Check all subdirectories and the episode file for a suitable series name
		episodePath := episode.Path

		// get the diff between episode.Path and episode.EpisodeFile so that we can
		// add one path element a time and extract the series
		subPath := episode.EpisodeFile[len(episodePath):]

		splits := strings.Split(subPath, "/")

		for _, part := range splits {
			if part == "" {
				continue
			}

			episodePath = path.Join(episodePath, part)

			subEpisode, subErr := renamer.CreateEpisodeFromPath(episodePath)
			if subErr == nil {
				possibilities = append(possibilities, subEpisode.Series)
			}
		}
	}

	return possibilities, nil
}

type ScriptExtractor struct {
	ScriptPath string
}

func (s ScriptExtractor) Names(episode *renamer.Episode) ([]string, error) {
	script := fmt.Sprintf("%s \"%s\" %d_%d",
		s.ScriptPath, episode.EpisodeFile, episode.Season, episode.Episode)
	cmd := exec.Command("/bin/sh", "-c", script)

	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(output), "\n"), nil
}

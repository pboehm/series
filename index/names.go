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

func (self FilesystemExtractor) Names(epi *renamer.Episode) ([]string, error) {
	possibilities := []string{epi.Series}

	if util.IsDirectory(epi.Path) {
		// Check all subdirectories and the episodefile itepi for a suitable
		// series name
		episodepath := epi.Path

		// get the diff between epi.Path and epi.Episodefile so that we can
		// add one path element a time and extract the series
		subpath := epi.Episodefile[len(episodepath):]

		splits := strings.Split(subpath, "/")

		for _, part := range splits {
			if part == "" {
				continue
			}

			episodepath = path.Join(episodepath, part)

			subepisode, suberr := renamer.CreateEpisodeFromPath(episodepath)
			if suberr == nil {
				possibilities = append(possibilities, subepisode.Series)
			}
		}
	}

	return possibilities, nil
}

type ScriptExtractor struct {
	ScriptPath string
}

func (self ScriptExtractor) Names(epi *renamer.Episode) ([]string, error) {
	script := fmt.Sprintf("%s \"%s\" %d_%d",
		self.ScriptPath, epi.Episodefile, epi.Season, epi.Episode)
	cmd := exec.Command("/bin/sh", "-c", script)

	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(output), "\n"), nil
}

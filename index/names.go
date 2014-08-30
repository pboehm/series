package index

import (
	"github.com/pboehm/series/renamer"
	"github.com/pboehm/series/util"
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

package main

import (
	"github.com/pboehm/series/renamer"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
)

var renameEpisodes, addToIndex bool

var renameAndIndexCmd = &cobra.Command{
	Use:   "rename_and_index",
	Short: "Renames and indexes the supplied episodes.",
	Run:   renameAndIndexHandler,
}

func renameAndIndexHandler(cmd *cobra.Command, args []string) {
	dir := appConfig.EpisodeDirectory
	if customEpisodeDirectory != "" {
		dir = customEpisodeDirectory
	}
	HandleError(os.Chdir(dir))

	interestingEntries := GetInterestingDirEntries()
	if len(interestingEntries) == 0 {
		os.Exit(0)
	}

	callPreProcessingHook()
	loadIndex()

	LOG.Println("### Process all interesting files ...")
	renameableEpisodes := HandleInterestingEpisodes(interestingEntries)

	if len(renameableEpisodes) > 0 && renameEpisodes {
		writeIndex()

		LOG.Println("### Renaming episodes ...")

		for _, episode := range renameableEpisodes {
			LOG.Printf("> %s: %s", episode.Series, episode.CleanedFileName())

			HandleError(episode.Rename("."))
			LOG.Printf("  [OK]\n")

			callEpisodeHook(episode.CleanedFileName(), episode.Series)
		}

		callPostProcessingHook()
	}
}

func GetInterestingDirEntries() []string {
	content, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	validRegex := regexp.MustCompile("^S\\d+E\\d+.-.\\w+.*\\.\\w+$")

	var interesting []string
	for _, entry := range content {
		entryPath := entry.Name()

		if !renamer.IsInterestingDirEntry(entryPath) {
			continue
		}
		if validRegex.Match([]byte(entryPath)) {
			continue
		}

		interesting = append(interesting, entryPath)
	}

	return interesting
}

func HandleInterestingEpisodes(entries []string) []*renamer.Episode {
	var renameableEpisodes []*renamer.Episode

	for _, entryPath := range entries {

		episode, err := renamer.CreateEpisodeFromPath(entryPath)
		if err != nil {
			LOG.Printf("!!! '%s' - %s\n\n", entryPath, err)
			continue
		}

		episode.RemoveTrashWords()
		if !episode.HasValidEpisodeName() {
			episode.SetDefaultEpisodeName()
		}

		LOG.Printf("<<< %s\n", entryPath)
		LOG.Printf(">>> %s\n", episode.CleanedFileName())

		if !episode.CanBeRenamed() {
			LOG.Printf("!!! '%s' is currently not renameable\n\n", entryPath)
			continue
		}

		if addToIndex {
			added, addedErr := seriesIndex.AddEpisode(episode)
			if !added {
				LOG.Printf("!!! couldn't be added to the index: %s\n\n", addedErr)
				continue
			}
			LOG.Printf("---> succesfully added to series index\n\n")
		}

		renameableEpisodes = append(renameableEpisodes, episode)
	}

	return renameableEpisodes
}

func init() {
	renameAndIndexCmd.Flags().BoolVarP(&renameEpisodes, "rename", "r", true,
		"Do actually rename the episodes.")
	renameAndIndexCmd.Flags().BoolVarP(&addToIndex, "index", "i", true,
		"Add the episodes to index.")

	indexCmd.Flags().BoolVarP(&renameEpisodes, "rename", "r", true,
		"Do actually rename the episodes.")
	indexCmd.Flags().BoolVarP(&addToIndex, "index", "i", true,
		"Add the episodes to index.")
}

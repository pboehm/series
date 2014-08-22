package main

import (
	"fmt"
	"github.com/pboehm/series/renamer"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
)

var RenameEpisodes, AddToIndex bool

var renameAndIndexCmd = &cobra.Command{
	Use:   "rename_and_index",
	Short: "Renames and indexes the supplied episodes.",
	Run:   renameAndIndexHandler,
}

func renameAndIndexHandler(cmd *cobra.Command, args []string) {
	dir := AppConfig.EpisodeDirectory
	if CustomEpisodeDirectory != "" {
		dir = CustomEpisodeDirectory
	}
	HandleError(os.Chdir(dir))

	interesting_entries := GetInterestingDirEntries()
	if len(interesting_entries) == 0 {
		os.Exit(0)
	}

	callPreProcessingHook()
	loadIndex()

	fmt.Println("### Process all interesting files ...")
	renameable_episodes := HandleInterestingEpisodes(interesting_entries)

	if len(renameable_episodes) > 0 && RenameEpisodes {
		writeIndex()

		fmt.Println("### Renaming episodes ...")

		for _, episode := range renameable_episodes {
			fmt.Printf("> %s: %s", episode.Series, episode.CleanedFileName())

			HandleError(episode.Rename("."))
			fmt.Printf("  [OK]\n")

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

	valid_regex := regexp.MustCompile("^S\\d+E\\d+.-.\\w+.*\\.\\w+$")

	interesting := []string{}
	for _, entry := range content {
		entry_path := entry.Name()

		if !renamer.IsInterestingDirEntry(entry_path) {
			continue
		}
		if valid_regex.Match([]byte(entry_path)) {
			continue
		}

		interesting = append(interesting, entry_path)
	}

	return interesting
}

func HandleInterestingEpisodes(entries []string) []*renamer.Episode {
	renameable_episodes := []*renamer.Episode{}

	for _, entry_path := range entries {

		episode, err := renamer.CreateEpisodeFromPath(entry_path)
		if err != nil {
			fmt.Printf("!!! '%s' - %s\n\n", entry_path, err)
			continue
		}

		episode.RemoveTrashwords()
		if !episode.HasValidEpisodeName() {
			episode.SetDefaultEpisodeName()
		}

		fmt.Printf("<<< %s\n", entry_path)
		fmt.Printf(">>> %s\n", episode.CleanedFileName())

		if !episode.CanBeRenamed() {
			fmt.Printf("!!! '%s' is currently not renameable\n\n", entry_path)
			continue
		}

		if AddToIndex {
			added, added_err := SeriesIndex.AddEpisode(episode)
			if !added {
				fmt.Printf("!!! couldn't be added to the index: %s\n\n", added_err)
				continue
			}
			fmt.Println("---> succesfully added to series index\n")
		}

		renameable_episodes = append(renameable_episodes, episode)
	}

	return renameable_episodes
}

func init() {
	renameAndIndexCmd.Flags().BoolVarP(&RenameEpisodes, "rename", "r", true,
		"Do actually rename the episodes.")
	renameAndIndexCmd.Flags().BoolVarP(&AddToIndex, "index", "i", true,
		"Add the episodes to index.")

	indexCmd.Flags().BoolVarP(&RenameEpisodes, "rename", "r", true,
		"Do actually rename the episodes.")
	indexCmd.Flags().BoolVarP(&AddToIndex, "index", "i", true,
		"Add the episodes to index.")
}

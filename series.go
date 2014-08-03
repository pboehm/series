package main

import (
	"flag"
	"fmt"
	"github.com/pboehm/series/config"
	"github.com/pboehm/series/index"
	"github.com/pboehm/series/renamer"
	"github.com/pboehm/series/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
)

var CONFIG_DIR, CONFIG_FILE string
var DEFAULT_CONFIG, APP_CONFIG config.Config

func setup() {
	CONFIG_DIR = path.Join(util.HomeDirectory(), ".series")
	CONFIG_FILE = path.Join(CONFIG_DIR, "config.json")

	DEFAULT_CONFIG = config.Config{
		EpisodeDirectory: path.Join(util.HomeDirectory(), "Downloads"),
		IndexFile:        path.Join(CONFIG_DIR, "index.xml"),
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

func HandleInterestingEpisodes(index *index.SeriesIndex, entries []string) []*renamer.Episode {
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

		added, added_err := index.AddEpisode(episode)
		if !added {
			fmt.Printf("!!! couldn't be added to the index: %s\n\n", added_err)
			continue
		}
		fmt.Println("---> succesfully added to series index\n")

		renameable_episodes = append(renameable_episodes, episode)
	}

	return renameable_episodes
}

// This executes the supplied cmd by /bin/sh and returns an error if it returns
// unexpectedly
func System(cmd_string string) error {

	cmd := exec.Command("/bin/sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	setup()
	APP_CONFIG = config.GetConfig(CONFIG_FILE, DEFAULT_CONFIG)

	// parse command flags/args
	FlagRenameFiles := flag.Bool("rename", true, "should the files be renamed")

	flag.Parse()
	argv := flag.Args()

	// change to the series directory
	dir := path.Join(APP_CONFIG.EpisodeDirectory)
	if len(argv) > 0 {
		dir = argv[0]
	}

	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	// get all interesting episodes and stop if there aren't any
	interesting_entries := GetInterestingDirEntries()
	if len(interesting_entries) == 0 {
		os.Exit(0)
	}

	// Call PreProcessingHook
	if APP_CONFIG.PreProcessingHook != "" {
		fmt.Println("### Calling PreProcessingHook ...")

		err := System(APP_CONFIG.PreProcessingHook)
		if err != nil {
			fmt.Printf("PreProcessingHook ended with an error: %s\n", err)
		}
	}

	fmt.Println("### Parsing series index ...")
	index, index_err := index.ParseSeriesIndex(APP_CONFIG.IndexFile)
	if index_err != nil {
		panic(index_err)
	}

	fmt.Println("### Process all interesting files ...")
	renameable_episodes := HandleInterestingEpisodes(index, interesting_entries)

	if len(renameable_episodes) > 0 && *FlagRenameFiles {
		fmt.Println("### Writing new index version ...")
		index.WriteToFile(APP_CONFIG.IndexFile)

		fmt.Println("### Renaming episodes ...")

		for _, episode := range renameable_episodes {
			fmt.Printf("> %s: %s", episode.Series, episode.CleanedFileName())

			// Rename episode file
			rename_err := episode.Rename(".")
			if rename_err != nil {
				panic(rename_err)
			}

			fmt.Printf("  [OK]\n")

			// Call EpisodeHook
			if APP_CONFIG.EpisodeHook != "" {
				fmt.Println("# Calling EpisodeHook ...")

				hook_cmd := fmt.Sprintf("%s \"%s\" \"%s\"",
					APP_CONFIG.EpisodeHook,
					episode.CleanedFileName(), episode.Series)

				err := System(hook_cmd)
				if err != nil {
					fmt.Printf("EpisodeHook ended with an error: %s\n", err)
				}
			}
		}

		// Call PostProcessingHook
		if APP_CONFIG.PostProcessingHook != "" {
			fmt.Println("\n### Calling PostProcessingHook ...")

			err := System(APP_CONFIG.PostProcessingHook)
			if err != nil {
				fmt.Printf("PostProcessingHook ended with an error: %s\n", err)
			}
		}
	}
}

package main

import (
	"errors"
	"fmt"
	"github.com/pboehm/series/index"
	"github.com/pboehm/series/util"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
)

var seriesIndex *index.SeriesIndex
var newSeriesLanguage string
var newSeriesFirstEpisode string

func loadIndex() {
	indexFilePath := appConfig.IndexFile
	if !util.PathExists(indexFilePath) {
		HandleError(errors.New(fmt.Sprintf("series index file %s does not exist, create one via `series index init`", indexFilePath)))
	}

	LOG.Println("### Parsing series index ...")

	var err error
	seriesIndex, err = index.ParseSeriesIndex(indexFilePath)
	HandleError(err)

	// add each SeriesNameExtractor
	seriesIndex.AddExtractor(index.FilesystemExtractor{})

	for _, script := range appConfig.ScriptExtractors {
		seriesIndex.AddExtractor(index.ScriptExtractor{ScriptPath: script})
	}
}

func writeIndex() {
	LOG.Println("### Writing new index version ...")
	seriesIndex.WriteToFile(appConfig.IndexFile)
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Manage the series index",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var indexInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an empty series index",
	Run: func(cmd *cobra.Command, args []string) {
		indexFilePath := appConfig.IndexFile
		if util.PathExists(indexFilePath) {
			HandleError(errors.New("series index already initialized"))
		}

		seriesIndex = &index.SeriesIndex{}

		writeIndex()
		callPostProcessingHook()
	},
}

var indexAddCmd = &cobra.Command{
	Use:   "add [series, ...]",
	Short: "Add series to index",
	Run: func(cmd *cobra.Command, args []string) {
		pattern := regexp.MustCompile("^S(?P<season>\\d+)E(?P<episode>\\d+)$")
		groups, matched := util.NamedCaptureGroups(pattern, newSeriesFirstEpisode)
		if !matched {
			HandleError(errors.New("first episode does not have the correct format like: S01E01"))
		}

		season, seasonErr := strconv.Atoi(groups["season"])
		episode, episodeErr := strconv.Atoi(groups["episode"])
		if seasonErr != nil || episodeErr != nil || season <= 0 || episode <= 0 {
			HandleError(errors.New("first episode is invalid"))
		}

		callPreProcessingHook()
		loadIndex()

		for _, seriesName := range args {
			LOG.Printf("Creating new index entry for '%s' [%s] with %s as first episode\n",
				seriesName, newSeriesLanguage, newSeriesFirstEpisode)

			_, err := seriesIndex.AddSeries(seriesName, newSeriesLanguage, season, episode-1)
			if err != nil {
				LOG.Printf(
					"!!! Adding new index entry wasn't possible: %s\n", err)
			}
		}

		writeIndex()
		callPostProcessingHook()
	},
}

var indexRemoveCmd = &cobra.Command{
	Use:   "remove [series, ...]",
	Short: "Remove series from index",
	Run: func(cmd *cobra.Command, args []string) {
		callPreProcessingHook()
		loadIndex()

		for _, seriesName := range args {
			LOG.Printf("Removing '%s' from index\n", seriesName)

			_, err := seriesIndex.RemoveSeries(seriesName)
			if err != nil {
				LOG.Printf("!!! Removing series from index wasn't possible: %s\n", err)
			}
		}

		writeIndex()
		callPostProcessingHook()
	},
}

var indexAliasCmd = &cobra.Command{
	Use:   "alias series [alias, ...]",
	Short: "Aliases the given series to the supplied aliases",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			LOG.Println("You have to supply one series name and some aliases")
			cmd.Usage()
			os.Exit(1)
		}

		callPreProcessingHook()
		loadIndex()

		series, args := args[0], args[1:]

		for _, alias := range args {
			LOG.Printf("Aliasing '%s' to '%s'\n", series, alias)
			err := seriesIndex.AliasSeries(series, alias)
			if err != nil {
				LOG.Printf("!!! Unable to alias the series: %s\n", err)
			}
		}

		writeIndex()
		callPostProcessingHook()
	},
}

var indexListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all series in index",
	Run: func(cmd *cobra.Command, args []string) {
		callPreProcessingHook()
		loadIndex()

		for _, series := range seriesIndex.SeriesList {
			fmt.Println(series.Name)
		}
	},
}

func init() {
	indexAddCmd.Flags().StringVarP(&newSeriesLanguage, "lang", "l", "de",
		"language the series is watched in. (de/en/fr)")
	indexAddCmd.Flags().StringVarP(&newSeriesFirstEpisode, "first-episode", "f", "S01E01",
		"the first episode that you are interested in")

	indexCmd.AddCommand(indexInitCmd, indexAddCmd, indexRemoveCmd, indexAliasCmd, indexListCmd)
}

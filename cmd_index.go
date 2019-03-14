package main

import (
	"fmt"
	"github.com/pboehm/series/index"
	"github.com/spf13/cobra"
	"os"
)

var seriesIndex *index.SeriesIndex
var newSeriesLanguage string

func loadIndex() {
	LOG.Println("### Parsing series index ...")

	var err error
	seriesIndex, err = index.ParseSeriesIndex(appConfig.IndexFile)
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

var addIndexCmd = &cobra.Command{
	Use:   "add [series, ...]",
	Short: "Add series to index",
	Run: func(cmd *cobra.Command, args []string) {
		callPreProcessingHook()
		loadIndex()

		for _, seriesName := range args {
			LOG.Printf("Creating new index entry for '%s' [%s]\n",
				seriesName, newSeriesLanguage)

			_, err := seriesIndex.AddSeries(seriesName, newSeriesLanguage)
			if err != nil {
				LOG.Printf(
					"!!! Adding new index entry wasn't possible: %s\n", err)
			}
		}

		writeIndex()
		callPostProcessingHook()
	},
}

var removeIndexCmd = &cobra.Command{
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

var aliasIndexCmd = &cobra.Command{
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

var listIndexCmd = &cobra.Command{
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
	addIndexCmd.Flags().StringVarP(&newSeriesLanguage, "lang", "l", "de",
		"language the series is watched in. (de/en/fr)")

	indexCmd.AddCommand(addIndexCmd, removeIndexCmd, aliasIndexCmd, listIndexCmd)
}

package main

import (
	"fmt"
	"github.com/pboehm/series/index"
	"github.com/spf13/cobra"
	"os"
)

var SeriesIndex *index.SeriesIndex
var NewSeriesLanguage string

func loadIndex() {
	fmt.Println("### Parsing series index ...")

	var err error
	SeriesIndex, err = index.ParseSeriesIndex(AppConfig.IndexFile)
	HandleError(err)
}

func writeIndex() {
	fmt.Println("### Writing new index version ...")
	SeriesIndex.WriteToFile(AppConfig.IndexFile)
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

		for _, seriesname := range args {
			fmt.Printf("Creating new index entry for '%s' [%s]\n",
				seriesname, NewSeriesLanguage)

			_, err := SeriesIndex.AddSeries(seriesname, NewSeriesLanguage)
			if err != nil {
				fmt.Printf(
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

		for _, seriesname := range args {
			fmt.Printf("Removing '%s' from index\n", seriesname)

			_, err := SeriesIndex.RemoveSeries(seriesname)
			if err != nil {
				fmt.Printf(
					"!!! Removing series from index wasn't possible: %s\n", err)
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
			fmt.Println("You have to supply one series name and some aliases")
			cmd.Usage()
			os.Exit(1)
		}

		callPreProcessingHook()
		loadIndex()

		series, args := args[0], args[1:]

		for _, alias := range args {
			fmt.Printf("Aliasing '%s' to '%s'\n", series, alias)
			err := SeriesIndex.AliasSeries(series, alias)
			if err != nil {
				fmt.Printf("!!! Unable to alias the series: %s\n", err)
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

		for _, series := range SeriesIndex.SeriesList {
			fmt.Printf("%s\n", series.Name)
		}
	},
}

func init() {
	addIndexCmd.Flags().StringVarP(&NewSeriesLanguage, "lang", "l", "de",
		"language the series is watched in. (de/en/fr)")
}

package main

import (
	"fmt"
	idx "github.com/pboehm/series/index"
	str "github.com/pboehm/series/streams"
	"github.com/spf13/cobra"
)

var streamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "Manage str that are interesting according to the index",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var streamsUnknownSeriesCmd = &cobra.Command{
	Use:   "unknown",
	Short: "List all series which are unknown by the streaming site",
	Run: func(cmd *cobra.Command, args []string) {
		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []str.WatchedSeries) {
			existingSeries := map[string]idx.Series{}
			for _, series := range index.SeriesList {
				existingSeries[series.Name] = series
			}

			for _, series := range watched {
				delete(existingSeries, series.SeriesNameInIndex)
			}

			for _, series := range existingSeries {
				fmt.Printf("Series is unknown by streaming site: %s\n", series.Name)
			}
		})
	},
}

var streamsFetchLinksCmd = &cobra.Command{
	Use:   "links",
	Short: "Fetch Links for unwatched episodes of series",
	Run: func(cmd *cobra.Command, args []string) {
		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []str.WatchedSeries) {
			linkSet := str.NewLinkSet(appConfig, streams, index)
			linkSet.GrabLinksFor(watched)
			grouped := linkSet.GroupedEntries()

			for group, groupedEntries := range grouped {
				fmt.Printf(">> Episodes for %s\n", group)
				for _, entry := range groupedEntries {
					link := entry.Links[0]
					fmt.Printf("%s - %s  [%s]\n", entry.Id, link.Link, link.Hoster)
				}
			}
		})
	},
}

func withIndexStreamsAndWatchedSeries(handler func(*idx.SeriesIndex, *str.Streams, []str.WatchedSeries)) {
	callPreProcessingHook()
	loadIndex()

	streams := &str.Streams{Config: appConfig}
	availableSeries := streams.AvailableSeries()

	var watched []str.WatchedSeries
	for _, series := range availableSeries {
		nameInIndex := seriesIndex.SeriesNameInIndex(series.Name)
		if nameInIndex != "" {
			languages := seriesIndex.SeriesLanguages(nameInIndex)
			watched = append(watched, str.WatchedSeries{
				Series:            series,
				SeriesNameInIndex: nameInIndex,
				SeriesLanguages:   mapLanguagesToIds(languages),
			})
		}
	}

	handler(seriesIndex, streams, watched)
}

func mapLanguagesToIds(languages []string) map[string]int {
	mapped := map[string]int{}

	for _, language := range languages {
		switch language {
		case "de":
			mapped[language] = 1
		case "en":
			mapped[language] = 2
		default:
			continue
		}
	}

	return mapped
}

func init() {
	streamsCmd.AddCommand(streamsUnknownSeriesCmd, streamsFetchLinksCmd)
}

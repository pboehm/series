package main

import (
	"fmt"
	idx "github.com/pboehm/series/index"
	str "github.com/pboehm/series/streams"
	"github.com/spf13/cobra"
)

type WatchedSeries struct {
	series            *str.Series
	seriesNameInIndex string
	seriesLanguages   map[string]int
}

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
		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []WatchedSeries) {
			existingSeries := map[string]idx.Series{}
			for _, series := range index.SeriesList {
				existingSeries[series.Name] = series
			}

			for _, series := range watched {
				delete(existingSeries, series.seriesNameInIndex)
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
		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []WatchedSeries) {
			for _, series := range watched {
				seasons := streams.Seasons(series.series)
				episodes := streams.Episodes(series.series, seasons[0])

				for _, episode := range episodes {
					for language, languageInt := range series.seriesLanguages {
						links := episode.LinksInLanguage(languageInt)
						if len(links) == 0 {
							continue
						}

						if !index.IsEpisodeInIndexManual(series.seriesNameInIndex, language, episode.Season, episode.Episode) {
							notifyNewEpisode(streams, series, language, episode, links)
						}
					}
				}
			}
		})
	},
}

func notifyNewEpisode(streams *str.Streams, series WatchedSeries, language string, episode *str.Episode, links []*str.Link) {
	fmt.Printf("S%02dE%02d from %s [%s] is available\n", episode.Season, episode.Episode, series.series.Name, language)
	for _, link := range links {
		fmt.Printf("  %15s  %s\n", link.Hoster, streams.LinkUrl(link))
	}
}

func withIndexStreamsAndWatchedSeries(handler func(*idx.SeriesIndex, *str.Streams, []WatchedSeries)) {
	callPreProcessingHook()
	loadIndex()

	streams := &str.Streams{Config: appConfig}
	availableSeries := streams.AvailableSeries()

	var watched []WatchedSeries
	for _, series := range availableSeries {
		nameInIndex := seriesIndex.SeriesNameInIndex(series.Name)
		if nameInIndex != "" {
			languages := seriesIndex.SeriesLanguages(nameInIndex)
			watched = append(watched, WatchedSeries{
				series:            series,
				seriesNameInIndex: nameInIndex,
				seriesLanguages:   mapLanguagesToIds(languages),
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

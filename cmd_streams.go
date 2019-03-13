package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	idx "github.com/pboehm/series/index"
	str "github.com/pboehm/series/streams"
	"github.com/spf13/cobra"
)

var streamsCmdJsonOutput = false

var streamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "Manage str that are interesting according to the index",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var streamsUnknownSeriesCmd = &cobra.Command{
	Use:   "unknown",
	Short: "list all series which are unknown by the streaming site",
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
				fmt.Println(series.Name)
			}
		})
	},
}

var streamsFetchLinksCmd = &cobra.Command{
	Use:   "links",
	Short: "fetch Links for unwatched episodes of series",
	Run: func(cmd *cobra.Command, args []string) {
		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []str.WatchedSeries) {
			linkSet := str.NewLinkSet(appConfig, streams, index)
			linkSet.GrabLinksFor(watched)

			if streamsCmdJsonOutput {
				entries := linkSet.Entries()
				bytes, err := json.MarshalIndent(entries, "", "  ")
				HandleError(err)
				fmt.Println(string(bytes))
			} else {
				grouped := linkSet.GroupedEntries()
				for group, groupedEntries := range grouped {
					fmt.Printf(">>>> %s\n", group)
					for _, entry := range groupedEntries {
						fmt.Printf(">> %s [%s]\n", entry.Filename, entry.Id)

						for i, link := range entry.Links {
							if i >= 2 {
								break
							}
							fmt.Printf("  %s\t  [%s]\n", link.Link, link.Hoster)
						}
					}
				}
			}
		})
	},
}

var streamsMarkWatchedCmd = &cobra.Command{
	Use:   "mark-watched [id, ....]",
	Short: "mark links as watched",
	Run: func(cmd *cobra.Command, args []string) {
		callPreProcessingHook()
		loadIndex()

		for _, arg := range args {
			id, e := str.IdentifierFromString(arg)
			HandleError(e)

			filename := fmt.Sprintf("S%02dE%02d - Episode %d.mov", id.Season, id.Episode, id.Episode)
			_, err := seriesIndex.AddEpisodeManually(id.Series, id.Language, id.Season, id.Episode, filename)
			if err == nil {
				LOG.Printf("Marking %s of %s [%s] as watched\n", filename, id.Series, id.Language)
			} else {
				LOG.Printf("Could not mark %s of %s [%s] as watched: %s\n", filename, id.Series, id.Language, err)
			}
		}

		writeIndex()
		callPostProcessingHook()
	},
}

var streamsServerCmd = &cobra.Command{
	Use:   "server",
	Short: "run an HTTP server serving an API and a frontend for streams",
	Run: func(cmd *cobra.Command, args []string) {
		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []str.WatchedSeries) {
			linkSet := str.NewLinkSet(appConfig, streams, index)
			linkSet.GrabLinksFor(watched)

			r := gin.Default()
			r.SetHTMLTemplate(str.ServerStaticHtmlTemplate)
			r.GET("/", func(c *gin.Context) {
				c.HTML(200, "index", gin.H{})
			})
			r.GET("/api/links", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"links": linkSet.Entries(),
				})
			})
			r.GET("/api/links/grouped", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"links": linkSet.GroupedEntries(),
				})
			})
			HandleError(r.Run())
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
	streamsFetchLinksCmd.Flags().BoolVarP(&streamsCmdJsonOutput, "json", "j", false, "output as JSON")

	streamsCmd.AddCommand(streamsUnknownSeriesCmd, streamsFetchLinksCmd, streamsMarkWatchedCmd, streamsServerCmd)
}

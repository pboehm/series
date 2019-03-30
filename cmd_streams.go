package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pboehm/series/config"
	idx "github.com/pboehm/series/index"
	str "github.com/pboehm/series/streams"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
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

func markEpisodeAsWatched(index *idx.SeriesIndex, episodeId string) (*str.Identifier, error) {
	var err error

	id, err := str.IdentifierFromString(episodeId)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("S%02dE%02d - Episode %d.mov", id.Season, id.Episode, id.Episode)
	_, err = index.AddEpisodeManually(id.Series, id.Language, id.Season, id.Episode, filename)
	return id, err
}

var streamsMarkWatchedCmd = &cobra.Command{
	Use:   "mark-watched [id, ....]",
	Short: "mark links as watched",
	Run: func(cmd *cobra.Command, args []string) {
		callPreProcessingHook()
		loadIndex()

		for _, arg := range args {
			id, err := markEpisodeAsWatched(seriesIndex, arg)
			if id == nil && err != nil {
				HandleError(err)
			}

			seasonWithEpisode := fmt.Sprintf("S%02dE%02d", id.Season, id.Episode)
			if err == nil {
				LOG.Printf("Marking %s of %s [%s] as watched\n", seasonWithEpisode, id.Series, id.Language)
			} else {
				LOG.Printf("Could not mark %s of %s [%s] as watched: %s\n", seasonWithEpisode, id.Series, id.Language, err)
			}
		}

		writeIndex()
		callPostProcessingHook()
	},
}

var streamsRunActionCmd = &cobra.Command{
	Use:   "run-action action linkId",
	Short: "run link actions on the supplied link ids",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Help()
			return
		}

		actionId := args[0]
		linkIdString := args[1]

		var action = config.StreamAction{}
		for _, a := range appConfig.StreamsLinkActions {
			if a.Id == actionId {
				action = a
			}
		}

		if action.Id == "" {
			HandleError(errors.New(fmt.Sprintf("action '%s' not found", actionId)))
		}

		withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []str.WatchedSeries) {
			identifier, linkId, err := str.LinkIdentifierFromString(linkIdString)
			HandleError(err)

			HandleError(runAction(os.Stderr, streams, action, identifier, linkId))
		})
	},
}

func runAction(output io.Writer, streams *str.Streams, action config.StreamAction, id *str.Identifier, linkId int) error {
	var err error
	var session, videoUrl string

	if session, err = streams.Login(appConfig.StreamsAccountEmail, appConfig.StreamsAccountPassword); err != nil {
		return err
	}

	if videoUrl, err = streams.ResolveLink(linkId, session); err != nil {
		return err
	}

	return SystemV(action.Command, []string{
		fmt.Sprintf("SERIES_SERIES=%s", id.Series),
		fmt.Sprintf("SERIES_LANGUAGE=%s", id.Language),
		fmt.Sprintf("SERIES_SEASON=%d", id.Season),
		fmt.Sprintf("SERIES_EPISODE=%d", id.Episode),
		fmt.Sprintf("SERIES_EPISODE_NAME=%s", id.EpisodeName),
		fmt.Sprintf("SERIES_FILENAME=S%02dE%02d - %s.mov", id.Season, id.Episode, id.EpisodeName),
		fmt.Sprintf("SERIES_LINK_ID=%d", linkId),
		fmt.Sprintf("SERIES_REDIRECT_URL=%s", streams.LinkUrl(linkId)),
		fmt.Sprintf("SERIES_VIDEO_URL=%s", videoUrl),
	}, output, output)
}

var streamsServerOptionListen string
var streamsServerOptionIndexHtml string

var streamsServerCmd = &cobra.Command{
	Use:   "server",
	Short: "run an HTTP server serving an API and a frontend for streams",
	Run: func(cmd *cobra.Command, args []string) {
		indexHtmlContent := func() []byte {
			if streamsServerOptionIndexHtml != "" {
				indexFileBytes, err := ioutil.ReadFile(streamsServerOptionIndexHtml)
				if err != nil {
					HandleError(err)
				}

				return indexFileBytes
			} else {
				return []byte(str.ServerStaticHtml)
			}
		}

		var currentLinkSet *str.LinkSet
		var currentStreams *str.Streams

		loadLinkSet := func() {
			withIndexStreamsAndWatchedSeries(func(index *idx.SeriesIndex, streams *str.Streams, watched []str.WatchedSeries) {
				linkSet := str.NewLinkSet(appConfig, streams, index)
				linkSet.GrabLinksFor(watched)
				currentLinkSet = linkSet
				currentStreams = streams
			})
		}

		go loadLinkSet()

		api := str.API{
			Config:         appConfig,
			Jobs:           str.NewJobPool(),
			HtmlContent:    indexHtmlContent,
			LinkSetRefresh: loadLinkSet,
			LinkSet: func() *str.LinkSet {
				return currentLinkSet
			},
			MarkWatched: func(episodeIds []string) ([]string, []string) {
				var successes, failures []string

				callPreProcessingHook()
				loadIndex()

				for _, episodeId := range episodeIds {
					_, err := markEpisodeAsWatched(seriesIndex, episodeId)
					if err == nil {
						successes = append(successes, episodeId)
					} else {
						failures = append(failures, episodeId)
					}
				}

				writeIndex()
				callPostProcessingHook()

				loadLinkSet()

				return successes, failures
			},
			ExecuteLinkAction: func(action config.StreamAction, identifier *str.Identifier, i int) *str.Job {
				return str.NewJob(func(output io.Writer) error {
					multiWriter := io.MultiWriter(output, os.Stderr)

					if currentStreams == nil {
						return errors.New("streams not initialized")
					}

					return runAction(multiWriter, currentStreams, action, identifier, i)
				})
			},
			ExecuteGlobalAction: func(action config.StreamAction) *str.Job {
				return str.NewJob(func(output io.Writer) error {
					multiWriter := io.MultiWriter(output, os.Stderr)
					return SystemV(action.Command, []string{}, multiWriter, multiWriter)
				})
			},
		}
		HandleError(api.Run(streamsServerOptionListen))
	},
}

func withIndexStreamsAndWatchedSeries(handler func(*idx.SeriesIndex, *str.Streams, []str.WatchedSeries)) {
	if appConfig.StreamsAPIToken == "" {
		HandleError(errors.New(fmt.Sprintf("`StreamsAPIToken` not configured in %s", configFile)))
	}

	if appConfig.StreamsAccountEmail == "" || appConfig.StreamsAccountPassword == "" {
		HandleError(errors.New(fmt.Sprintf(
			"`StreamsAccountEmail` or `StreamsAccountPassword` is not configured in %s", configFile)))
	}

	callPreProcessingHook()
	loadIndex()

	streams := str.NewStreams(appConfig)
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

	streamsServerCmd.Flags().StringVarP(&streamsServerOptionListen, "listen", "l", ":8080", "where should the server listen")
	streamsServerCmd.Flags().StringVarP(&streamsServerOptionIndexHtml, "index", "i", "", "a custom index.html that should be used")

	streamsCmd.AddCommand(streamsUnknownSeriesCmd, streamsFetchLinksCmd, streamsMarkWatchedCmd, streamsRunActionCmd, streamsServerCmd)
}

package streams

import (
	"fmt"
	"github.com/pboehm/series/config"
	"github.com/pboehm/series/index"
	"sort"
	"strings"
)

type LinkSetEntryLink struct {
	Hoster string `json:"hoster"`
	Link   string `json:"link"`
}

type LinkSetEntry struct {
	Id       string              `json:"id"`
	Series   string              `json:"series"`
	Language string              `json:"language"`
	Season   int                 `json:"season"`
	Episode  int                 `json:"episode"`
	Links    []*LinkSetEntryLink `json:"links"`
}

type LinkSet struct {
	config       config.Config
	streams      *Streams
	index        *index.SeriesIndex
	episodeLinks []*LinkSetEntry
}

func NewLinkSet(config config.Config, streams *Streams, index *index.SeriesIndex) *LinkSet {
	return &LinkSet{
		config:  config,
		streams: streams,
		index:   index,
	}
}

func (l *LinkSet) GrabLinksFor(watched []WatchedSeries) {
	for _, series := range watched {
		seasons := l.streams.Seasons(series.Series)
		episodes := l.streams.Episodes(series.Series, seasons[0])

		for _, episode := range episodes {
			for language, languageInt := range series.SeriesLanguages {
				links := episode.LinksInLanguage(languageInt)
				if len(links) == 0 {
					continue
				}

				if !l.index.IsEpisodeInIndexManual(series.SeriesNameInIndex, language, episode.Season, episode.Episode) {
					l.notifyNewEpisode(series, language, episode, links)
				}
			}
		}
	}
}

func (l *LinkSet) notifyNewEpisode(series WatchedSeries, language string, episode *Episode, links []*Link) {
	var entryLinks []*LinkSetEntryLink
	for _, link := range links {
		entryLinks = append(entryLinks, &LinkSetEntryLink{
			Hoster: link.Hoster,
			Link:   l.streams.LinkUrl(link),
		})
	}

	episodeLink := &LinkSetEntry{
		Id:       fmt.Sprintf("S%02dE%02d", episode.Season, episode.Episode),
		Series:   series.SeriesNameInIndex,
		Language: language,
		Season:   episode.Season,
		Episode:  episode.Episode,
		Links:    entryLinks,
	}

	l.episodeLinks = append(l.episodeLinks, episodeLink)
}

func (l *LinkSet) Entries() []*LinkSetEntry {
	return l.episodeLinks
}

func (l *LinkSet) GroupedEntries() map[string][]*LinkSetEntry {
	grouped := map[string][]*LinkSetEntry{}
	for _, entry := range l.Entries() {
		groupId := fmt.Sprintf("%s [%s]", entry.Series, entry.Language)

		entriesByGroupId, ok := grouped[groupId]
		if !ok {
			entriesByGroupId = []*LinkSetEntry{}
		}

		// TODO replace by real hoster selection
		sort.Slice(entry.Links[:], func(i, j int) bool {
			return strings.Index(entry.Links[i].Hoster, "HD") > strings.Index(entry.Links[j].Hoster, "HD")
		})

		grouped[groupId] = append(entriesByGroupId, entry)
	}

	return grouped
}


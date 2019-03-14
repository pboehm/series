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
	Id          string              `json:"id"`
	Series      string              `json:"series"`
	Language    string              `json:"language"`
	Season      int                 `json:"season"`
	Episode     int                 `json:"episode"`
	EpisodeId   int                 `json:"episode_id"`
	EpisodeName string              `json:"episode_name"`
	Filename    string              `json:"filename"`
	Links       []*LinkSetEntryLink `json:"links"`
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
					l.addEpisode(series, language, episode, links)
				}
			}
		}
	}

	sort.Slice(l.episodeLinks, func(i, j int) bool {
		return l.episodeLinks[i].EpisodeId < l.episodeLinks[j].EpisodeId
	})
}

func (l *LinkSet) addEpisode(series WatchedSeries, language string, episode *Episode, links []*Link) {
	var entryLinks []*LinkSetEntryLink
	for _, link := range links {
		entryLinks = append(entryLinks, &LinkSetEntryLink{
			Hoster: link.Hoster,
			Link:   l.streams.LinkUrl(link),
		})
	}

	// TODO replace by real hoster selection
	sort.Slice(entryLinks[:], func(i, j int) bool {
		return strings.Index(entryLinks[i].Hoster, "HD") > strings.Index(entryLinks[j].Hoster, "HD")
	})

	id, _ := NewIdentifier(series.SeriesNameInIndex, language, episode.Season, episode.Episode).AsString()

	episodeName := ""
	switch language {
	case "de":
		episodeName = episode.German
	case "en":
		episodeName = episode.English
	default:
	}

	episodeLink := &LinkSetEntry{
		Id:          id,
		Series:      series.SeriesNameInIndex,
		Language:    language,
		Season:      episode.Season,
		Episode:     episode.Episode,
		EpisodeId:   episode.ID,
		EpisodeName: episodeName,
		Filename:    fmt.Sprintf("S%02dE%02d - %s.mov", episode.Season, episode.Episode, episodeName),
		Links:       entryLinks,
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

		grouped[groupId] = append(entriesByGroupId, entry)
	}

	return grouped
}

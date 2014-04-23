package index

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/pboehm/series/renamer"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var DefaultLanguage = "de"

type SeriesIndex struct {
	XMLName    xml.Name `xml:"seriesindex"`
	SeriesList []Series `xml:"series"`
	seriesMap  map[string]*Series
}

func (self *SeriesIndex) AddEpisode(episode *renamer.Episode) (bool, error) {

    // test for all possible series names if they exist in index and take the
    // first matching
    for _, possible_series := range episode.GetPossibleSeriesNames() {
        series_name := self.SeriesNameInIndex(possible_series)
        if series_name != "" {
            episode.Series = series_name
            break
        }
    }

	if episode.Series == "" {
		return false, errors.New("Series does not exist in index")
	}
	series := self.seriesMap[episode.Series]

	// Handle episodes where no language is set
	if episode.Language == "" {
		self.GuessEpisodeLanguage(episode, series)
	}

	_, language_exist := series.languageMap[episode.Language]
	if !language_exist {
		return false, errors.New("Series is not watched in this language")
	}

	if self.IsEpisodeInIndex(*episode) {
		return false, errors.New("Episode already exists in Index")
	}

	episode_entry := Episode{Name: episode.CleanedFileName()}

	// find the right EpisodeSet so we can add our new episode to it
	set, exist := series.languageMap[episode.Language]
	if exist {
		set.EpisodeList = append(set.EpisodeList, episode_entry)

		// add it to the lookup cache
		key := GetIndexKey(episode.Season, episode.Episode)
		set.episodeMap[key] = episode_entry.Name

		return true, nil
	}

	return false,
		errors.New("Episode couldn't be added to index. This shouldn't occur!")
}

func (self *SeriesIndex) GuessEpisodeLanguage(episode *renamer.Episode, series *Series) {
	// This methods tries to find the right language for the supplied episode
	// based on several heuristics

	// When there is no language set and the series is only watched in one
	// language we can take this series
	if len(series.languageMap) == 1 {
		for k, _ := range series.languageMap {
			episode.Language = k
			break
		}
	}

	// Find the language which is most likely the right language
	if len(series.languageMap) > 1 {
		possible_languages := []string{}

		// when episode has not been watched in only one of the languages
		for lang, _ := range series.languageMap {
			episode.Language = lang
			if !self.IsEpisodeInIndex(*episode) {
				possible_languages = append(possible_languages, lang)
			}

			episode.Language = ""
		}

		if len(possible_languages) == 1 {
			episode.Language = possible_languages[0]

		} else if len(possible_languages) > 1 {
			// take the language where the previous episode exist
			previous_existing := []string{}

			for _, lang := range possible_languages {
				epi := *episode
				epi.Language = lang
				if (epi.Episode - 1) > 0 {
					epi.Episode -= 1
				}

				if self.IsEpisodeInIndex(epi) {
					previous_existing = append(previous_existing, lang)
				}
			}

			if len(previous_existing) == 1 {
				episode.Language = previous_existing[0]
			}
		}
	}
}

func (self *SeriesIndex) SeriesNameInIndex(series_name string) string {

	series_in_index, exist := self.seriesMap[series_name]
	if exist {
		return series_in_index.Name
	}

	// do a case insensitive search
	joined := series_name
	for {
		if joined == "" {
			break
		}

		pattern := regexp.MustCompile(fmt.Sprintf("^(?i)%s$", joined))
		for name, series := range self.seriesMap {
			if pattern.Match([]byte(name)) {
				return series.Name
			}
		}

		splitted := strings.Split(joined, " ")
		joined = strings.Join(splitted[1:], " ")
	}

	return ""
}

func (self *SeriesIndex) IsEpisodeInIndex(episode renamer.Episode) bool {

	series_name := self.SeriesNameInIndex(episode.Series)
	if series_name == "" {
		return false
	}

	series, series_exist := self.seriesMap[series_name]
	if !series_exist {
		return false
	}

	set, language_exist := series.languageMap[episode.Language]
	if !language_exist {
		return false
	}

	key := GetIndexKey(episode.Season, episode.Episode)
	_, episode_exist := set.episodeMap[key]

	if episode_exist {
		return true
	}

	// check if episode is before the lowest episode which sets all_before=true
	// takes place
	if set.allBefore {
		barrier := set.allBeforeSeason*100 + set.allBeforeEpisode
		actual := episode.Season*100 + episode.Episode

		if actual < barrier {
			return true
		}
	}

	return false
}

func ParseSeriesIndex(xmlpath string) (*SeriesIndex, error) {
	var index SeriesIndex

	xmlFile, err := os.Open(xmlpath)
	if err != nil {
		return &index, err
	}
	defer xmlFile.Close()

	content, err := ioutil.ReadAll(xmlFile)

	xml.Unmarshal([]byte(content), &index)

	index.BuildUpSeriesMap()
	return &index, nil
}

func (index *SeriesIndex) BuildUpSeriesMap() {
	// Build up the series map that holds references to series under the series
	// name and all aliases
	index.seriesMap = map[string]*Series{}

	for i := 0; i < len(index.SeriesList); i++ {
		series := &(index.SeriesList[i])
		series.BuildUpLanguageMap()

		index.seriesMap[series.Name] = series

		for _, alias := range series.Aliases {
			index.seriesMap[alias.To] = series
		}
	}
}

func (self *SeriesIndex) WriteToFile(xmlpath string) {

	marshaled, err := xml.MarshalIndent(*self, "", "  ")
	if err != nil {
		panic(err)
	}

	output := append([]byte(xml.Header), marshaled...)

	err = ioutil.WriteFile(xmlpath, output, 0644)
	if err != nil {
		panic(err)
	}
}

type Series struct {
	Name        string       `xml:"name,attr"`
	EpisodeSets []EpisodeSet `xml:"episodes"`
	Aliases     []Alias      `xml:"alias"`
	languageMap map[string]*EpisodeSet
}

func (self *Series) BuildUpLanguageMap() {
	self.languageMap = make(map[string]*EpisodeSet)

	for i := 0; i < len(self.EpisodeSets); i++ {
		set := &(self.EpisodeSets[i])
		set.BuildUpEpisodeMap()
		self.languageMap[set.GetLanguage()] = set
	}
}

type EpisodeSet struct {
	XMLName                           xml.Name  `xml:"episodes"`
	EpisodeList                       []Episode `xml:"episode"`
	Language                          string    `xml:"lang,attr,omitempty"`
	episodeMap                        map[string]string
	allBefore                         bool
	allBeforeSeason, allBeforeEpisode int
}

func (self *EpisodeSet) BuildUpEpisodeMap() {
	self.episodeMap = make(map[string]string)

	for _, episode := range self.EpisodeList {

		matched := renamer.ExtractEpisodeInformation(episode.Name)
		if matched != nil {
			nr_season, _ := strconv.Atoi(matched["season"])
			nr_episode, _ := strconv.Atoi(matched["episode"])
			key := GetIndexKey(nr_season, nr_episode)

			self.episodeMap[key] = episode.Name

			// handle all_before flag and set data for later usage
			if episode.AllBefore {
				self.allBefore = true
				self.allBeforeSeason = nr_season
				self.allBeforeEpisode = nr_episode
			}
		}
	}
}

func (self *EpisodeSet) GetLanguage() string {
	if self.Language != "" {
		return self.Language
	}

	return DefaultLanguage
}

type Episode struct {
	Name      string `xml:"name,attr"`
	AllBefore bool   `xml:"all_before,attr,omitempty"`
}

type Alias struct {
	To string `xml:"to,attr"`
}

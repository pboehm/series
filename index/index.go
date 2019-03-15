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
	XMLName        xml.Name `xml:"seriesindex"`
	SeriesList     []Series `xml:"series"`
	seriesMap      map[string]*Series
	nameExtractors []SeriesNameExtractor
}

// AddExtractor adds another SeriesNameExtractor for generating possible series
// names. Extractors get called in order order of addition.
func (s *SeriesIndex) AddExtractor(ex SeriesNameExtractor) {
	s.nameExtractors = append(s.nameExtractors, ex)
}

func (s *SeriesIndex) AddEpisode(episode *renamer.Episode) (bool, error) {

	// test for all possible series names if they exist in index and take the
	// first matching
ExtractorLoop:
	for _, extractor := range s.nameExtractors {
		names, err := extractor.Names(episode)

		if err != nil {
			fmt.Printf("!!! Error asking extractor for series names: %s", err)
			continue
		}

		for _, possibleSeries := range names {
			seriesName := s.SeriesNameInIndex(possibleSeries)
			if seriesName != "" {
				episode.Series = seriesName
				break ExtractorLoop
			}
		}
	}

	series, existing := s.seriesMap[episode.Series]
	if !existing {
		return false, errors.New("series does not exist in index")
	}

	// Handle episodes where no language is set
	if episode.Language == "" {
		s.GuessEpisodeLanguage(episode, series)
	}

	return s.AddEpisodeManually(episode.Series, episode.Language, episode.Season, episode.Episode, episode.CleanedFileName())
}

func (s *SeriesIndex) AddEpisodeManually(seriesNameInIndex string, language string, season int, episode int, filename string) (bool, error) {
	series, existing := s.seriesMap[seriesNameInIndex]
	if !existing {
		return false, errors.New("series does not exist in index")
	}

	_, languageExist := series.languageMap[language]
	if !languageExist {
		return false, errors.New("series is not watched in this language")
	}

	if s.IsEpisodeInIndexManual(series.Name, language, season, episode) {
		return false, errors.New("episode already exists in index")
	}

	episodeEntry := Episode{Name: filename}

	// find the right EpisodeSet so we can add our new episode to it
	set, exist := series.languageMap[language]
	if exist {
		set.EpisodeList = append(set.EpisodeList, episodeEntry)
		set.BuildUpEpisodeMap()
		return true, nil
	}

	return false,
		errors.New("episode couldn't be added to index, this shouldn't occur")
}

func (s *SeriesIndex) AddSeries(seriesname, language string, season int, episode int) (bool, error) {

	_, existing := s.seriesMap[seriesname]
	if existing {
		return false, errors.New("series does already exist in index")
	}

	series := Series{
		Name: seriesname,
		EpisodeSets: []EpisodeSet{
			{
				Language: language,
				EpisodeList: []Episode{
					{
						Name:      fmt.Sprintf("S%02dE%02d - Pre-First.mov", season, episode),
						AllBefore: true,
					},
				},
			},
		},
	}

	s.SeriesList = append(s.SeriesList, series)
	s.BuildUpSeriesMap()

	return true, nil
}

func (s *SeriesIndex) RemoveSeries(seriesname string) (bool, error) {

	series, existing := s.seriesMap[seriesname]
	if !existing {
		return false, errors.New("series does not exist in index")
	}

	for i := 0; i < len(s.SeriesList); i++ {
		if s.SeriesList[i].Name == series.Name {
			s.SeriesList = append(
				s.SeriesList[:i],
				s.SeriesList[i+1:]...,
			)
			s.BuildUpSeriesMap()
			break
		}
	}

	return true, nil
}

func (s *SeriesIndex) AliasSeries(seriesname string, alias string) error {

	series, existing := s.seriesMap[seriesname]
	if !existing {
		return errors.New("series does not exist in index")
	}

	_, aliasExisting := s.seriesMap[alias]
	if aliasExisting {
		return errors.New("alias does already exist as series in index")
	}

	series.Aliases = append(series.Aliases, Alias{To: alias})
	series.BuildUpLanguageMap()

	return nil
}

func (s *SeriesIndex) GuessEpisodeLanguage(episode *renamer.Episode, series *Series) {
	// This methods tries to find the right language for the supplied episode
	// based on several heuristics

	// When there is no language set and the series is only watched in one
	// language we can take this series
	if len(series.languageMap) == 1 {
		for k := range series.languageMap {
			episode.Language = k
			break
		}
	}

	// Find the language which is most likely the right language
	if len(series.languageMap) > 1 {
		var possibleLanguages []string

		// when episode has not been watched in only one of the languages
		for lang := range series.languageMap {
			episode.Language = lang
			if !s.IsEpisodeInIndex(*episode) {
				possibleLanguages = append(possibleLanguages, lang)
			}

			episode.Language = ""
		}

		if len(possibleLanguages) == 1 {
			episode.Language = possibleLanguages[0]

		} else if len(possibleLanguages) > 1 {
			// take the language where the previous episode exist
			var previousExisting []string

			for _, lang := range possibleLanguages {
				epi := *episode
				epi.Language = lang
				if (epi.Episode - 1) > 0 {
					epi.Episode -= 1
				}

				if s.IsEpisodeInIndex(epi) {
					previousExisting = append(previousExisting, lang)
				}
			}

			if len(previousExisting) == 1 {
				episode.Language = previousExisting[0]
			}
		}
	}
}

func (s *SeriesIndex) SeriesNameInIndex(seriesName string) string {

	seriesInIndex, exist := s.seriesMap[seriesName]
	if exist {
		return seriesInIndex.Name
	}

	// do a case insensitive search
	joined := seriesName
	for {
		if joined == "" {
			break
		}

		pattern := regexp.MustCompile(fmt.Sprintf("^(?i)%s$", regexp.QuoteMeta(joined)))
		for name, series := range s.seriesMap {
			if pattern.Match([]byte(name)) {
				return series.Name
			}
		}

		splitted := strings.Split(joined, " ")
		joined = strings.Join(splitted[1:], " ")
	}

	return ""
}

func (s *SeriesIndex) IsEpisodeInIndex(episode renamer.Episode) bool {
	return s.IsEpisodeInIndexManual(episode.Series, episode.Language, episode.Season, episode.Episode)
}

func (s *SeriesIndex) IsEpisodeInIndexManual(inputSeriesName string, language string, season int, episode int) bool {

	seriesName := s.SeriesNameInIndex(inputSeriesName)
	if seriesName == "" {
		return false
	}

	series, seriesExist := s.seriesMap[seriesName]
	if !seriesExist {
		return false
	}

	set, languageExist := series.languageMap[language]
	if !languageExist {
		return false
	}

	key := buildIndexKey(season, episode)
	_, episodeExist := set.episodeMap[key]

	if episodeExist {
		return true
	}

	// check if episode is before the lowest episode which sets all_before=true
	// takes place
	if set.allBefore {
		barrier := set.allBeforeSeason*100 + set.allBeforeEpisode
		actual := season*100 + episode

		if actual < barrier {
			return true
		}
	}

	return false
}

func (s *SeriesIndex) SeriesLanguages(seriesNameInIndex string) []string {
	var languages []string

	series, ok := s.seriesMap[seriesNameInIndex]
	if ok {
		for key, _ := range series.languageMap {
			languages = append(languages, key)
		}
	}

	return languages
}

func ParseSeriesIndex(xmlPath string) (*SeriesIndex, error) {
	var index SeriesIndex

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return &index, err
	}
	defer xmlFile.Close()

	content, err := ioutil.ReadAll(xmlFile)

	xml.Unmarshal([]byte(content), &index)

	index.BuildUpSeriesMap()
	return &index, nil
}

func (s *SeriesIndex) BuildUpSeriesMap() {
	// Build up the series map that holds references to series under the series
	// name and all aliases
	s.seriesMap = map[string]*Series{}

	for i := 0; i < len(s.SeriesList); i++ {
		series := &(s.SeriesList[i])
		series.BuildUpLanguageMap()

		s.seriesMap[series.Name] = series

		for _, alias := range series.Aliases {
			s.seriesMap[alias.To] = series
		}
	}
}

func (s *SeriesIndex) WriteToFile(xmlPath string) {

	marshaled, err := xml.MarshalIndent(*s, "", "  ")
	if err != nil {
		panic(err)
	}

	output := append([]byte(xml.Header), marshaled...)

	err = ioutil.WriteFile(xmlPath, output, 0644)
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

func (s *Series) BuildUpLanguageMap() {
	s.languageMap = make(map[string]*EpisodeSet)

	for i := 0; i < len(s.EpisodeSets); i++ {
		set := &(s.EpisodeSets[i])
		set.BuildUpEpisodeMap()
		s.languageMap[set.GetLanguage()] = set
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

func (e *EpisodeSet) BuildUpEpisodeMap() {
	e.episodeMap = make(map[string]string)

	for _, episode := range e.EpisodeList {

		matched := renamer.ExtractEpisodeInformation(episode.Name)
		if matched != nil {
			nrSeason, _ := strconv.Atoi(matched["season"])
			nrEpisode, _ := strconv.Atoi(matched["episode"])
			key := buildIndexKey(nrSeason, nrEpisode)

			e.episodeMap[key] = episode.Name

			// handle all_before flag and set data for later usage
			if episode.AllBefore {
				e.allBefore = true
				e.allBeforeSeason = nrSeason
				e.allBeforeEpisode = nrEpisode
			}
		}
	}
}

func (e *EpisodeSet) GetLanguage() string {
	if e.Language != "" {
		return e.Language
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

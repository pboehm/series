package index

import (
	"encoding/xml"
	"fmt"
	"github.com/pboehm/series/renamer"
	"io/ioutil"
	"os"
	"strconv"
	"regexp"
	"strings"
	"errors"
)

var DefaultLanguage = "de"

type SeriesIndex struct {
	XMLName    xml.Name `xml:"seriesindex"`
	SeriesList []Series `xml:"series"`
	seriesMap  map[string]*Series
}

func (self *SeriesIndex) AddEpisode(episode *renamer.Episode) (bool, error) {
    series_name := self.SeriesNameInIndex(episode.Series)
    if series_name == "" {
        return false, errors.New("Series does not exist in index")
    }
    episode.Series = series_name
    series := self.seriesMap[episode.Series]


    // Handle episodes where no language is set
    if episode.Language == "" {

        // When there is no language set and the series is only watched in one
        // language we can take this series
        if len(series.episodeMap) == 1 {
            for k, _ := range series.episodeMap {
                episode.Language = k
                break
            }
        }

        if len(series.episodeMap) > 1 {
            // Find the language which is most likely the right language
            // - when episode exists in one of the languages
            // - take the language where the episode is the nearest one
        }
    }

    _, language_exist := series.episodeMap[episode.Language]
    if ! language_exist {
        return false, errors.New("Series is not watched in this language")
    }

    if self.IsEpisodeInIndex(*episode) {
        return false, errors.New("Episode already exists in Index")
    }

    episode_entry := Episode{Name: episode.CleanedFileName()}

    // find the right EpisodeSet so we can add our new episode to it
    for i := 0; i < len(series.EpisodeSets); i++ {
        set := &(series.EpisodeSets[i])
        if set.GetLanguage() == episode.Language {
            set.EpisodeList = append(set.EpisodeList, episode_entry)

            // add it to the lookup cache
            key := GetIndexKey(episode.Season, episode.Episode)
            series.episodeMap[episode.Language][key] = episode_entry.Name

            return true, nil
        }
    }

    return false,
        errors.New("Episode couldn't be added to index. This should no occur!")
}


func (self *SeriesIndex) SeriesNameInIndex(series_name string) string {

    series_in_index, exist := self.seriesMap[series_name]
    if exist {
        return series_in_index.Name
    }

    // do a case insensitive search
    joined := series_name
    for {
        if (joined == "") { break }

        pattern := regexp.MustCompile(fmt.Sprintf("^(?i)%s$", joined))
        for name, series := range self.seriesMap {
            if pattern.Match([]byte(name)) {
                return series.Name
            }
        }

        splitted := strings.Split(joined, " ")
        joined = strings.Join(splitted[1:], " ")
    }

    return "";
}

func (self *SeriesIndex) IsEpisodeInIndex(episode renamer.Episode) bool {

    series_name := self.SeriesNameInIndex(episode.Series)
    if series_name == "" { return false }

    series, series_exist := self.seriesMap[series_name]
    if ! series_exist { return false }

    _, language_exist := series.episodeMap[episode.Language]
    if ! language_exist { return false }

    key := GetIndexKey(episode.Season, episode.Episode)
    _, episode_exist := series.episodeMap[episode.Language][key]

    return episode_exist
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

	index.SetupLookupCaches()
	return &index, nil
}

func (index *SeriesIndex) SetupLookupCaches() {
	// Build up the series map that holds references to series under the series
	// name and all aliases
	index.seriesMap = map[string]*Series{}

	for i := 0; i < len(index.SeriesList); i++ {
		series := &(index.SeriesList[i])
		series.BuildUpEpisodeMap()

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

func (self *SeriesIndex) Print() {
	index := *self

	fmt.Println(index.SeriesList)
	for _, series := range index.SeriesList {
		fmt.Printf(">>> %s\n", series.Name)
		for _, episodeset := range series.EpisodeSets {
			fmt.Printf(">>>> %s - %d\n", episodeset.Language,
				len(episodeset.EpisodeList))
			for _, episode := range episodeset.EpisodeList {
				fmt.Printf(">>>>> %s\n", episode.Name)
			}
		}
	}
}

type Series struct {
	Name        string       `xml:"name,attr"`
	EpisodeSets []EpisodeSet `xml:"episodes"`
	Aliases     []Alias      `xml:"alias"`
	episodeMap  map[string]map[string]string
}

func (self *Series) BuildUpEpisodeMap() {
	self.episodeMap = make(map[string]map[string]string)

	for _, set := range self.EpisodeSets {
		self.episodeMap[set.GetLanguage()] = make(map[string]string)

		for _, episode := range set.EpisodeList {
			matched := renamer.ExtractEpisodeInformation(episode.Name)
			if matched != nil {
				nr_season, _ := strconv.Atoi(matched["season"])
				nr_episode, _ := strconv.Atoi(matched["episode"])
                key := GetIndexKey(nr_season, nr_episode)

				self.episodeMap[set.GetLanguage()][key] = episode.Name
			}
		}
	}
}

type EpisodeSet struct {
	XMLName     xml.Name  `xml:"episodes"`
	EpisodeList []Episode `xml:"episode"`
	Language    string    `xml:"lang,attr,omitempty"`
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

package index

import (
	"encoding/xml"
	"fmt"
	"github.com/pboehm/series/renamer"
	"io/ioutil"
	"os"
	"strconv"
)

type SeriesIndex struct {
	XMLName    xml.Name `xml:"seriesindex"`
	SeriesList []Series `xml:"series"`
	SeriesMap  map[string]*Series
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
	index.SeriesMap = map[string]*Series{}

	for i := 0; i < len(index.SeriesList); i++ {
		series := &(index.SeriesList[i])
		series.BuildUpEpisodeMap()

		index.SeriesMap[series.Name] = series

		for _, alias := range series.Aliases {
			index.SeriesMap[alias.To] = series
		}
	}
}

func (self *SeriesIndex) WriteToFile(xmlpath string) error {
	fmt.Printf("%s", xmlpath)

	marshaled, err := xml.MarshalIndent(*self, "", "  ")
	fmt.Printf("%s", xml.Header)
	fmt.Printf("%s\n", marshaled)

	return err
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
	EpisodeMap  map[string]map[string]string
}

func (self *Series) BuildUpEpisodeMap() {
	self.EpisodeMap = make(map[string]map[string]string)

	for _, set := range self.EpisodeSets {
		self.EpisodeMap[set.GetLanguage()] = make(map[string]string)

		for _, episode := range set.EpisodeList {
			matched := renamer.ExtractEpisodeInformation(episode.Name)
			if matched != nil {
				nr_season, _ := strconv.Atoi(matched["season"])
				nr_episode, _ := strconv.Atoi(matched["episode"])
				key := fmt.Sprintf("%d_%d", nr_season, nr_episode)

				self.EpisodeMap[set.GetLanguage()][key] = episode.Name
			}
		}
	}
}

type EpisodeSet struct {
	XMLName     xml.Name  `xml:"episodes"`
	EpisodeList []Episode `xml:"episode"`
	Language    string    `xml:"lang,attr"`
}

func (self *EpisodeSet) GetLanguage() string {
	if self.Language != "" {
		return self.Language
	}

	return "de"
}

type Episode struct {
	Name      string `xml:"name,attr"`
	AllBefore bool   `xml:"all_before,attr,omitempty"`
}

type Alias struct {
	To string `xml:"to,attr"`
}

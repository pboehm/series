package index

import (
	"os"
	"encoding/xml"
	"io/ioutil"
	"fmt"
)

type SeriesIndex struct {
    XMLName xml.Name `xml:"seriesindex"`
    SeriesList []Series `xml:"series"`
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

	return &index, nil
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
    Name string `xml:"name,attr"`
    EpisodeSets []EpisodeSet `xml:"episodes"`
    Aliases []Alias `xml:"alias"`
}

type EpisodeSet struct {
    XMLName xml.Name `xml:"episodes"`
    EpisodeList []Episode `xml:"episode"`
    Language string `xml:"lang,attr"`
}

type Episode struct {
    Name string `xml:"name,attr"`
    AllBefore bool `xml:"all_before,attr,omitempty"`
}

type Alias struct {
    To string `xml:"to,attr"`
}

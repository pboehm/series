package streams

import (
	"encoding/base64"
	"encoding/json"
)

type WatchedSeries struct {
	Series            *Series
	SeriesNameInIndex string
	SeriesLanguages   map[string]int
}

type Identifier struct {
	Series   string `json:"series"`
	Language string `json:"language"`
	Season   int    `json:"season"`
	Episode  int    `json:"episode"`
}

func (i *Identifier) AsString() (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func NewIdentifier(series string, language string, season int, episode int) *Identifier {
	return &Identifier{Series: series, Language: language, Season: season, Episode: episode}
}

func IdentifierFromString(id string) (*Identifier, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}

	var identifier Identifier
	if err = json.Unmarshal(bytes, &identifier); err != nil {
		return nil, err
	}

	return &identifier, err
}

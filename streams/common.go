package streams

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type WatchedSeries struct {
	Series            *Series
	SeriesNameInIndex string
	SeriesLanguages   map[string]int
}

type Identifier struct {
	Series      string `json:"series"`
	Language    string `json:"language"`
	Season      int    `json:"season"`
	Episode     int    `json:"episode"`
	EpisodeName string `json:"episode_name"`
}

func (i *Identifier) AsString() (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func NewIdentifier(series string, language string, season int, episode int, name string) *Identifier {
	return &Identifier{Series: series, Language: language, Season: season, Episode: episode, EpisodeName: name}
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

func NewLinkIdentifier(identifier *Identifier, linkId int) (string, error) {
	idString, err := identifier.AsString()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s___%d", idString, linkId), nil
}

func LinkIdentifierFromString(idString string) (*Identifier, int, error) {
	parts := strings.Split(idString, "___")
	if len(parts) != 2 {
		return nil, -1, errors.New("invalid link identifier format")
	}

	linkId, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, -1, errors.New("invalid link id")
	}

	identifier, err := IdentifierFromString(parts[0])
	if err != nil {
		return nil, -1, errors.New("invalid identifier")
	}

	return identifier, linkId, nil
}

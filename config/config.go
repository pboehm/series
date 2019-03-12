package config

import (
	"encoding/json"
	"github.com/pboehm/series/util"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	IndexFile, PreProcessingHook, PostProcessingHook, EpisodeHook string
	EpisodeDirectory                                              string
	ScriptExtractors                                              []string
	StreamsAPIToken                                               string
}

func GetConfig(configFile string, standard Config) Config {

	if !util.PathExists(configFile) {
		configDir := path.Dir(configFile)

		dirErr := os.MkdirAll(configDir, 0755)
		if dirErr != nil {
			panic(dirErr)
		}

		writeMarshaledDataToFile(configFile, standard)
	}

	content, readErr := ioutil.ReadFile(configFile)
	if readErr != nil {
		panic(readErr)
	}

	unmarshalErr := json.Unmarshal(content, &standard)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}

	writeMarshaledDataToFile(configFile, standard)

	return standard
}

func writeMarshaledDataToFile(file string, config Config) {
	marshaled, marshalErr := json.MarshalIndent(config, "", "  ")
	if marshalErr != nil {
		panic(marshalErr)
	}

	writeErr := ioutil.WriteFile(file, marshaled, 0644)
	if writeErr != nil {
		panic(writeErr)
	}
}

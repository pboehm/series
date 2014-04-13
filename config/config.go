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
}

func GetConfig(config_file string, standard Config) Config {

	if !util.PathExists(config_file) {
		config_dir := path.Dir(config_file)

		dir_err := os.MkdirAll(config_dir, 0755)
		if dir_err != nil {
			panic(dir_err)
		}

        writeMarshaledDataToFile(config_file, standard)
	}

	content, read_err := ioutil.ReadFile(config_file)
    if read_err != nil {
        panic(read_err)
    }

	unmarshal_err := json.Unmarshal(content, &standard)
	if unmarshal_err != nil {
		panic(unmarshal_err)
	}

    writeMarshaledDataToFile(config_file, standard)

	return standard
}

func writeMarshaledDataToFile(file string, config Config) {
    marshaled, marshal_err := json.MarshalIndent(config, "", "  ")
    if marshal_err != nil {
        panic(marshal_err)
    }

    write_err := ioutil.WriteFile(file, marshaled, 0644)
    if write_err != nil {
        panic(write_err)
    }
}

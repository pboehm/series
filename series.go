package main

import (
	"github.com/pboehm/series/config"
	"github.com/pboehm/series/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

var LOG = log.New(os.Stderr, "", 0)

func HandleError(err error) {
	if err != nil {
		LOG.Fatalf("Error during execution: %s", err)
	}
}

var configDirectory, configFile, customEpisodeDirectory string
var defaultConfig, appConfig config.Config

func setupConfig() {
	configDirectory = path.Join(util.HomeDirectory(), ".series")
	configFile = path.Join(configDirectory, "config.json")

	defaultConfig = config.Config{
		EpisodeDirectory: path.Join(util.HomeDirectory(), "Downloads"),
		IndexFile:        path.Join(configDirectory, "index.xml"),
		ScriptExtractors: []string{},
	}

	appConfig = config.GetConfig(configFile, defaultConfig)
}

var seriesCmd = &cobra.Command{
	Use: "series",
	Run: renameAndIndexHandler,
}

func init() {
	seriesCmd.PersistentFlags().StringVarP(&customEpisodeDirectory, "dir", "d", "",
		"The directory which includes the episodes. (Overrides the config value)")
}

func main() {
	setupConfig()

	seriesCmd.AddCommand(renameAndIndexCmd, indexCmd, streamsCmd)
	seriesCmd.Execute()
}

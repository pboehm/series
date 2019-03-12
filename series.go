package main

import (
	"github.com/pboehm/series/config"
	"github.com/pboehm/series/util"
	"github.com/spf13/cobra"
	"log"
	"path"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
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

	seriesCmd.AddCommand(renameAndIndexCmd, indexCmd)
	indexCmd.AddCommand(addIndexCmd, removeIndexCmd, aliasIndexCmd, listIndexCmd)
	seriesCmd.Execute()
}

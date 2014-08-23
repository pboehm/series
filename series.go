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

var ConfigDirectory, ConfigFile, CustomEpisodeDirectory string
var DefaultConfig, AppConfig config.Config

func setupConfig() {
	ConfigDirectory = path.Join(util.HomeDirectory(), ".series")
	ConfigFile = path.Join(ConfigDirectory, "config.json")

	DefaultConfig = config.Config{
		EpisodeDirectory: path.Join(util.HomeDirectory(), "Downloads"),
		IndexFile:        path.Join(ConfigDirectory, "index.xml"),
	}

	AppConfig = config.GetConfig(ConfigFile, DefaultConfig)
}

var seriesCmd = &cobra.Command{
	Use: "series",
	Run: renameAndIndexHandler,
}

func init() {
	seriesCmd.PersistentFlags().StringVarP(&CustomEpisodeDirectory, "dir", "d", "",
		"The directory which includes the episodes. (Overrides the config value)")
}

func main() {
	setupConfig()

	seriesCmd.AddCommand(renameAndIndexCmd, indexCmd)
	indexCmd.AddCommand(addIndexCmd, removeIndexCmd, aliasIndexCmd, listIndexCmd)
	seriesCmd.Execute()
}

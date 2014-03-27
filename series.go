package main

import (
    "github.com/pboehm/series/renamer"
    "github.com/pboehm/series/util"
    "fmt"
    "io/ioutil"
    "path"
    "os"
    "regexp"
    "flag"
)

func main() {
    // parse command flags/args
    FlagRenameFiles := flag.Bool("rename", true, "should the files be renamed")

    flag.Parse()
    argv := flag.Args()

    // change to the series directory
    dir := path.Join(util.HomeDirectory(), "Downloads")
    if len(argv) > 0 {
        dir = argv[0]
    }

    err := os.Chdir(dir)
    if err != nil { panic(err) }

    content, err := ioutil.ReadDir(".")
    if err != nil { panic(err) }

    // handle all files that are interesting and not already in the right format
    valid_regex := regexp.MustCompile("^S\\d+E\\d+.-.\\w+.*\\.\\w+$")

    for _, entry := range content {
        entry_path := entry.Name()

        if ! renamer.IsInterestingDirEntry(entry_path) { continue }
        if valid_regex.Match([]byte(entry_path))       { continue }

        episode, err := renamer.CreateEpisodeFromPath(entry_path)
        if err != nil {
            fmt.Printf("!!! '%s' - %s\n\n", entry_path, err)
            continue
        }

        episode.RemoveTrashwords()
        if ! episode.HasValidEpisodeName() {
            episode.SetDefaultEpisodeName()
        }

        fmt.Printf("<<< %s\n", entry_path)
        fmt.Printf(">>> %s\n", episode.CleanedFileName())

        if ! episode.CanBeRenamed() {
            fmt.Printf("!!! %s couldn't be renamed\n", entry_path)
            continue
        }

        if *FlagRenameFiles {
            rename_err := episode.Rename(".")
            if rename_err != nil { panic(rename_err) }

            fmt.Println("--> episode has been renamed succesfully")
        }

        fmt.Println("")
    }
}

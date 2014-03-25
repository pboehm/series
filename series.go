package main

import (
    "github.com/pboehm/series/renamer"
    "github.com/pboehm/series/util"
    "fmt"
    "io/ioutil"
    "path"
    "os"
    "regexp"
)

func main() {
    dir := path.Join(util.HomeDirectory(), "Downloads")

    err := os.Chdir(dir)
    if err != nil { panic(err) }

    content, err := ioutil.ReadDir(".")
    if err != nil { panic(err) }

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

        rename_err := episode.Rename(".")
        if rename_err != nil { panic(rename_err) }

        fmt.Println("--> episode has been renamed succesfully\n")
    }
}

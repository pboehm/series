package renamer

import (
    "path"
    "io/ioutil"
    "regexp"
    "github.com/pboehm/series/util"
    "strings"
)

var Patterns = []*regexp.Regexp {
    // S01E01
    regexp.MustCompile(
        "^(?i)(?P<series>.*)S(?P<season>\\d+)E(?P<episode>\\d+)(?P<episodename>.*)$"),

    // 101; 1212
    regexp.MustCompile(
        "^(?i)(?P<series>.*\\D)(?P<season>\\d+)(?P<episode>\\d{2})(?P<episodename>\\W*.*)$"),

    // 1x1; 12x12
    regexp.MustCompile(
        "^(?i)(?P<series>.*)(?P<season>\\d+)x(?P<episode>\\d+)(?P<episodename>.*)$"),
}


func IsInterestingDirEntry(entry string) bool {
    for _, pattern := range Patterns {
        _, matched := util.NamedCaptureGroups(pattern, entry)
        if matched {
            return true
        }
    }
    return false
}

func ExtractEpisodeInformation(entry string) map[string]string {
    for _, pattern := range Patterns {
        groups, matched := util.NamedCaptureGroups(pattern, entry)
        if matched {
            return groups
        }
    }
    return nil
}

func CleanEpisodeInformation(info string) string {
    return strings.TrimSpace(strings.Replace(info, ".", " ", -1))
}

func GetDirtyFiles() []string {
    content, _ := ioutil.ReadDir(path.Join(util.HomeDirectory(), "Downloads"))

    var entries []string
    for _, entry := range content {
        if IsInterestingDirEntry(entry.Name()) {
            entries = append(entries, entry.Name())
        }
    }

    return entries
}


package renamer

import (
    "os"
    "path"
    "io/ioutil"
    "regexp"
)

var Patterns = []string {
    // S01E01
    "^(?i)(?P<series>.*)S(?P<season>\\d+)E(?P<episode>\\d+)(?P<episodename>.*)$",
    // 101; 1212
    "^(?i)(?P<series>.*\\D)(?P<season>\\d+)(?P<episode>\\d{2})(?P<episodename>\\W*.*)$",
    // 1x1; 12x12
    "^(?i)(?P<series>.*)(?P<season>\\d+)x(?P<episode>\\d+)(?P<episodename>.*)$",
}

func GetDirtyFiles() []string {
    content, _ := ioutil.ReadDir(path.Join(GetHomeDirectory(), "Downloads"))

    var entries []string
    for _, entry := range content {
        if IsInterestingDirEntry(entry.Name()) {
            entries = append(entries, entry.Name())
        }
    }

    return entries
}

func IsInterestingDirEntry(entry string) bool {
    for _, pattern := range Patterns {
        re := regexp.MustCompile(pattern)
        if re.Match([]byte(entry)) {
            return true
        }
    }
    return false
}

func GetHomeDirectory() string {
    return os.Getenv("HOME")
}

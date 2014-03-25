package renamer

import (
    "path"
    "regexp"
    "github.com/pboehm/series/util"
    "strings"
    "fmt"
    "errors"
    "path/filepath"
    "os"
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

var VideoFileEndings = []string {
    "mpg", "mpeg", "avi", "mkv", "wmv", "mp4", "mov", "flv", "3gp", "ts",
}

var TrashWords = []string {
    "German", "Dubbed", "DVDRip", "HDTVRip", "XviD", "ITG", "TVR", "inspired",
    "HDRip", "AMBiTiOUS", "RSG", "SiGHT", "SATRip", "WS", "TVS", "RiP", "READ",
    "GERMAN", "dTV", "aTV", "iNTERNAL", "CRoW", "MSE", "c0nFuSed", "UTOPiA",
    "scum", "EXPiRED", "BDRiP", "HDTV", "iTunesHD", "720p", "x264", "h264",
    "CRiSP", "euHD", "WEBRiP", "ZZGtv", "ARCHiV", "DD20", "Prim3time", "Nfo",
    "Repack", "SiMPTY", "BLURAYRiP", "BluRay", "DELiCiOUS", "Synced",
    "UNDELiCiOUS", "fBi", "CiD", "iTunesHDRip", "RedSeven", "OiNK", "idTV",
    "DL", "DD51", "AC3", "1080p",
}

func HasVideoFileEnding(entry_path string) (bool) {
    extension := path.Ext(entry_path)

    if extension == "" {
        return false
    } else {
        extension = extension[1:]
        matched := false

        for _, ending := range VideoFileEndings {
            if ending == extension {
                matched = true
                break
            }
        }

        return matched
    }
}

func FindBiggestVideoFile(dir string) (string, error) {
    if ! util.IsDirectory(dir) {
        return "", errors.New("The supplied directory does not exist")
    }

    var videofile string
    var videofile_size int64

    walker := func(entry_path string, info os.FileInfo, err error) error {
        if info.IsDir() || ! HasVideoFileEnding(entry_path) {
            return nil
        }

        if info.Size() > videofile_size {
            videofile = entry_path
            videofile_size = info.Size()
        }

        return nil
    }

    err := filepath.Walk(dir, walker)
    if err != nil {
        panic(err)
    }

    if videofile == "" {
        return "", errors.New("No videofile available")
    }

    return videofile, nil
}

func ApplyTrashwordsOnString(str string) string {
    purge_count := 0
    last_purge  := ""
    valid_words := []string {}
    splitted := regexp.MustCompile("\\s").Split(str, -1)

    SplittedLoop:
    for _, word := range splitted {
        if word == "" { continue SplittedLoop }

        word_pattern := regexp.MustCompile(fmt.Sprintf("^(?i)%s$", word))

        // Check if the current word is a known trashword
        for _, trashword := range TrashWords {
            if word_pattern.Match([]byte(trashword)) {
                purge_count++
                last_purge = word
                continue SplittedLoop
            }
        }

        // check if a valid word occurs after the first purged word
        if purge_count == 1 && last_purge != "" {
            valid_words = append(valid_words, last_purge)
            purge_count = 0
            last_purge = ""
        }

        valid_words = append(valid_words, word)
    }
    return strings.Join(valid_words, " ")
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

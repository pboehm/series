package renamer

import (
	"errors"
	"fmt"
	"github.com/pboehm/series/util"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	Patterns = []*regexp.Regexp{
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

	MultipleWhitespacePattern = regexp.MustCompile("\\s+")

	VideoFileEndings = []string{
		"mpg", "mpeg", "avi", "mkv", "wmv", "mp4", "mov", "flv", "3gp", "ts",
	}

	TrashWords = []string{
		"German", "Dubbed", "DVDRip", "HDTVRip", "XviD", "ITG", "TVR", "inspired",
		"HDRip", "AMBiTiOUS", "RSG", "SiGHT", "SATRip", "WS", "TVS", "RiP", "READ",
		"GERMAN", "dTV", "aTV", "iNTERNAL", "CRoW", "MSE", "c0nFuSed", "UTOPiA",
		"scum", "EXPiRED", "BDRiP", "HDTV", "iTunesHD", "720p", "x264", "h264",
		"CRiSP", "euHD", "WEBRiP", "ZZGtv", "ARCHiV", "DD20", "Prim3time", "Nfo",
		"Repack", "SiMPTY", "BLURAYRiP", "BluRay", "DELiCiOUS", "Synced",
		"UNDELiCiOUS", "fBi", "CiD", "iTunesHDRip", "RedSeven", "OiNK", "idTV",
		"DL", "DD51", "AC3", "1080p", "WEB", "DD5",
	}
)

func HasVideoFileEnding(entryPath string) bool {
	extension := path.Ext(entryPath)

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
	if !util.IsDirectory(dir) {
		return "", errors.New("the supplied directory does not exist")
	}

	var videoFile string
	var videoFileSize int64

	walker := func(entryPath string, info os.FileInfo, err error) error {
		if info.IsDir() || !HasVideoFileEnding(entryPath) {
			return nil
		}

		if info.Size() > videoFileSize {
			videoFile = entryPath
			videoFileSize = info.Size()
		}

		return nil
	}

	err := filepath.Walk(dir, walker)
	if err != nil {
		panic(err)
	}

	if videoFile == "" {
		return "", errors.New("no video file available")
	}

	return videoFile, nil
}

func ApplyTrashWordsOnString(str string) string {
	purgeCount := 0
	lastPurge := ""
	var validWords []string
	splitted := regexp.MustCompile("\\s").Split(str, -1)

SplittedLoop:
	for _, word := range splitted {
		if word == "" {
			continue SplittedLoop
		}
		if purgeCount > 2 {
			break
		}

		wordPattern := regexp.MustCompile(fmt.Sprintf("^(?i)%s$", word))

		// Check if the current word is a known trashWord
		for _, trashWord := range TrashWords {
			if wordPattern.Match([]byte(trashWord)) {
				purgeCount++
				lastPurge = word
				continue SplittedLoop
			}
		}

		// check if a valid word occurs after the first purged word
		if purgeCount == 1 && lastPurge != "" {
			validWords = append(validWords, lastPurge)
			purgeCount = 0
			lastPurge = ""
		}

		validWords = append(validWords, word)
	}
	return strings.Join(validWords, " ")
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
	cleaned := strings.Replace(
		strings.Replace(info, "-", " ", -1),
		".", " ", -1)

	cleaned = string(
		MultipleWhitespacePattern.ReplaceAll([]byte(cleaned), []byte(" ")))
	return strings.TrimSpace(cleaned)
}

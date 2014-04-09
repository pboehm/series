package util

import (
	"regexp"
)

// NamedCaptureGroups does a regexp match wih the supplied pattern on the str
// parameter and extracts the namend capture groups as map and returns it
//
// when the pattern does not match (nil, false) is returned
func NamedCaptureGroups(pattern *regexp.Regexp, str string) (map[string]string, bool) {
	match := pattern.FindStringSubmatch(str)
	if match == nil {
		return nil, false
	}

	result := make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		result[name] = match[i]
	}

	return result, true
}

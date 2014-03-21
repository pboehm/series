package util

import (
    "os"
)

// returns whether the given file or directory exists or not
func PathExists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    return false
}

func HomeDirectory() string {
    return os.Getenv("HOME")
}

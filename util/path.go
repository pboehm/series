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

func IsFile(path string) (bool) {
    stat, err := os.Stat(path)
    if err == nil && stat.Mode().IsRegular() {
        return true
    }
    return false
}

func IsDirectory(path string) (bool) {
    stat, err := os.Stat(path)
    if err == nil && stat.Mode().IsDir() {
        return true
    }
    return false
}

func HomeDirectory() string {
    return os.Getenv("HOME")
}

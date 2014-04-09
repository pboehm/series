package index

import (
    "fmt"
)

func GetIndexKey(season, episode int) string {
    return fmt.Sprintf("%d_%d", season, episode)
}

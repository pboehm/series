package index

import (
	"fmt"
)

func buildIndexKey(season, episode int) string {
	return fmt.Sprintf("%d_%d", season, episode)
}

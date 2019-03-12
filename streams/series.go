package streams

type WatchedSeries struct {
	Series            *Series
	SeriesNameInIndex string
	SeriesLanguages   map[string]int
}

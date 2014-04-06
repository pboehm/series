package index

import (
    . "launchpad.net/gocheck"
    "github.com/pboehm/series/util"
    "testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
var _ = Suite(&MySuite{})

type MySuite struct{
    dir string
    index *SeriesIndex
    err error
}

func (s *MySuite) SetUpTest(c *C) {
    s.dir = c.MkDir()
    s.index, s.err = ParseSeriesIndex("data/seriesindex_example.xml")
}

func (s *MySuite) TestEnvironment(c *C) {
    c.Assert(util.PathExists(s.dir), Equals, true)
}

func (s *MySuite) TestIndexParsing(c *C) {
    c.Assert(s.index, NotNil)
    c.Assert(s.err, IsNil)
    c.Assert(s.index.SeriesList, HasLen, 4)

    c.Assert(s.index.SeriesMap, HasLen, 6)
    c.Assert(s.index.SeriesMap["Community"], Equals, s.index.SeriesMap["Comm"])
    c.Assert(s.index.SeriesMap["Community"].Name, Equals, "Community")
}

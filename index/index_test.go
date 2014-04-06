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
}

func (s *MySuite) SetUpTest(c *C) {
    s.dir = c.MkDir()
}

func (s *MySuite) TestEnvironment(c *C) {
    c.Assert(util.PathExists(s.dir), Equals, true)
}

func (s *MySuite) TestIndexParsing(c *C) {
    index, err := ParseSeriesIndex("data/seriesindex_example.xml")

    c.Assert(index, NotNil)
    c.Assert(err, IsNil)
}

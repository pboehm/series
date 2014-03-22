package util

import (
    . "launchpad.net/gocheck"
    "testing"
    "regexp"
)

func Test(t *testing.T) { TestingT(t) }
var _ = Suite(&MySuite{})
type MySuite struct { }

func (s *MySuite) TestRegexpNamedCaptures(c *C) {
    pattern := regexp.MustCompile("^S(?P<season>\\d+)E(?P<episode>\\d+)$")

    groups, matched := NamedCaptureGroups(pattern, "S01E02")
    c.Assert(matched, Equals, true)
    c.Assert(groups, Not(IsNil))
    c.Assert(groups["season"], Equals, "01")
    c.Assert(groups["episode"], Equals, "02")

    groups, matched = NamedCaptureGroups(pattern, "SHOULDNOTMATCH")
    c.Assert(matched, Equals, false)
    c.Assert(groups, IsNil)
}

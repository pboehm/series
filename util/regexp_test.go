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
    c.Check(matched, Equals, true)
    c.Check(groups, Not(IsNil))
    c.Check(groups["season"], Equals, "01")
    c.Check(groups["episode"], Equals, "02")

    groups, matched = NamedCaptureGroups(pattern, "SHOULDNOTMATCH")
    c.Check(matched, Equals, false)
    c.Check(groups, IsNil)
}

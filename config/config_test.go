package config

import (
	"github.com/pboehm/series/util"
	. "launchpad.net/gocheck"
	"path"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&MySuite{})

type MySuite struct {
	dir        string
	configFile string
}

func (s *MySuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	s.configFile = path.Join(s.dir, ".series/config.json")
}

func (s *MySuite) TestConfigParsingWhenNoConfigExists(c *C) {
	standard := Config{}
	config := GetConfig(s.configFile, standard)

	c.Assert(config, DeepEquals, standard)
	c.Assert(util.PathExists(s.configFile), Equals, true)
}

func (s *MySuite) TestConfigParsingWithChanges(c *C) {
	// create config which seems to be changed by the user
	old := Config{IndexFile: "/not/existing/file.json"}
	_ = GetConfig(s.configFile, old)
	c.Assert(util.PathExists(s.configFile), Equals, true)

	standard := Config{IndexFile: "/other/non/existing/file"}
	config := GetConfig(s.configFile, standard)
	c.Assert(config.IndexFile, Equals, "/not/existing/file.json")
}

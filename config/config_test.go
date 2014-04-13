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
	dir         string
	config_file string
}

func (s *MySuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	s.config_file = path.Join(s.dir, ".series/config.json")
}

// type Config struct {
// 	IndexFile, PreProcessingHook, PostProcessingHook, EpisodeHook string
// 	EpisodeDirectory                                              string
// }

func (s *MySuite) TestConfigParsingWhenNoConfigExists(c *C) {
	standard := Config{}
	config := GetConfig(s.config_file, standard)

	c.Assert(config, Equals, standard)
	c.Assert(util.PathExists(s.config_file), Equals, true)
}

func (s *MySuite) TestConfigParsingWithChanges(c *C) {
	// create config which seems to be changed by the user
	old := Config{IndexFile: "/not/existing/file.json"}
	_ = GetConfig(s.config_file, old)
	c.Assert(util.PathExists(s.config_file), Equals, true)

	standard := Config{IndexFile: "/other/non/existing/file"}
	config := GetConfig(s.config_file, standard)
	c.Assert(config.IndexFile, Equals, "/not/existing/file.json")
}

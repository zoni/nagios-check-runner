package nca

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

// Config describes the full agent configuration.
type Config struct {
	Hostname string
	Checks   map[string]Check
}

// LoadConfig loads configuration from the given source and returns a
// fully initialized Configuration struct from it.
func LoadConfig(src io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	for name, check := range c.Checks {
		if check.Name == "" {
			check.Name = name
		}

		splitArgs, err := shellquote.Split(check.Command)
		if err != nil {
			return nil, err
		}
		if len(splitArgs) < 1 {
			return nil, fmt.Errorf("Check '%s' is missing a command to execute", name)
		}
		check.Args = splitArgs

		c.Checks[name] = check
	}

	return c, nil
}

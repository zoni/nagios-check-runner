package nca

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"strings"
)

// Config describes the full agent configuration.
type Config struct {
	Publishers map[string]map[string]interface{}
	Hostname   string
	Checks     map[string]Check
}

// ReadConfig loads configuration from the given source and returns a
// fully initialized Configuration struct from it.
func ReadConfig(src io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err = yaml.Unmarshal(data, c); err != nil {
		return nil, Error{
			Code:    ErrInvalidConfig,
			Message: err.Error(),
		}
	}

	if err = parseChecks(c); err != nil {
		return nil, err
	}

	if err = parsePublishers(c); err != nil {
		return nil, err
	}

	return c, nil
}

// parseChecks is a helper function to ReadConfig.
func parseChecks(cfg *Config) error {
	for name, check := range cfg.Checks {
		if check.Name == "" {
			check.Name = name
		}
		if check.Interval < 1 {
			check.Interval = 60
		}
		if check.Retry < 1 {
			check.Retry = 60
		}
		if check.Timeout < 1 {
			check.Timeout = 10
		}

		splitArgs, err := shellquote.Split(check.Command)
		if err != nil {
			return err
		}
		if len(splitArgs) < 1 {
			return Error{
				Code:    ErrInvalidConfig,
				Message: fmt.Sprintf("Check '%s' is missing a command to execute", name),
			}
		}
		check.Args = splitArgs

		cfg.Checks[name] = check
	}
	return nil
}

// parsePublishers is a helper function to ReadConfig.
func parsePublishers(cfg *Config) error {
	for label, publisher := range cfg.Publishers {
		_, found := publisher["type"]
		if !found {
			publisher["type"] = label + "publisher"
		}
		t, ok := publisher["type"].(string)
		if !ok {
			return Error{
				Code:    ErrInvalidConfig,
				Message: fmt.Sprintf("Type field of publisher '%s' should be a string", label),
			}
		}
		ptype := strings.ToLower(t)
		publisher["type"] = ptype

		if _, found := publisherFactories[ptype]; !found {
			return Error{
				Code:    ErrInvalidConfig,
				Message: fmt.Sprintf("No publisher named %q available", ptype),
			}
		}
	}

	return nil
}

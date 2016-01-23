package nca

import (
	"bytes"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var defaultConfig []byte

func init() {
	f, err := os.Open("testdata/config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	defaultConfig, err = ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
}

func TestEmptyConfig(t *testing.T) {
	a := assert.New(t)

	cfg, err := LoadConfig(strings.NewReader(""))
	require.Nil(t, err)

	a.Equal(0, len(cfg.Checks))
}

func TestRegularConfig(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	cfg, err := LoadConfig(bytes.NewReader(defaultConfig))
	r.Nil(err)

	a.Equal("testhost", cfg.Hostname)
	a.Equal(5, len(cfg.Checks))
}

func TestChecks(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	_, err := LoadConfig(strings.NewReader("checks: { valid: {command: foo} }"))
	a.Nil(err)
	_, err = LoadConfig(strings.NewReader("checks: { invalid: {} }"))
	a.NotNil(err)

	cfg, err := LoadConfig(bytes.NewReader(defaultConfig))
	r.Nil(err)

	check, ok := cfg.Checks["Dummy OK"]
	if !ok {
		t.Fatal("Expected 'Dummy OK' check not found")
	}

	a.Equal("Dummy OK", check.Name)
	a.Equal("/usr/lib/nagios/plugins/check_dummy 0", check.Command)
	a.Equal([]string{"/usr/lib/nagios/plugins/check_dummy", "0"}, check.Args)
	a.Equal(60, check.Interval)
	a.Equal(60, check.Retry)
	a.Equal(10, check.Timeout)

	check, ok = cfg.Checks["Custom"]
	if !ok {
		t.Fatal("Expected 'Dummy OK' check not found")
	}
	a.Equal("Custom check", check.Name)
	a.Equal(5, check.Interval)
	a.Equal(3, check.Retry)
	a.Equal(3, check.Timeout)
}

package nca

import (
	"bytes"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
	"io"
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

	cfg, err := ReadConfig(strings.NewReader(""))
	require.Nil(t, err)

	a.Equal(0, len(cfg.Checks))
	a.Equal(0, len(cfg.Publishers))
}

func TestRegularConfig(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	cfg, err := ReadConfig(bytes.NewReader(defaultConfig))
	r.Nil(err)

	a.Equal(5, len(cfg.Checks))
	a.Equal(1, len(cfg.Publishers))
}

func TestBadConfig(t *testing.T) {
	_, err := ReadConfig(strings.NewReader("invalid: {"))
	if err == nil {
		t.Error("Expected error not to be nil")
	}
	if errt, ok := err.(Error); !ok || errt.Code != ErrInvalidConfig {
		t.Error("Expected error to be ErrInvalidConfig")
	}
}

func TestChecks(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	_, err := ReadConfig(strings.NewReader("checks: { valid: {command: foo} }"))
	a.Nil(err)
	_, err = ReadConfig(strings.NewReader("checks: { invalid: {} }"))
	a.NotNil(err)
	if errt, ok := err.(Error); !ok || errt.Code != ErrInvalidConfig {
		t.Error("Expected error to be ErrInvalidConfig")
	}

	cfg, err := ReadConfig(bytes.NewReader(defaultConfig))
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
	a.Equal(10, check.Interval)
	a.Equal(3, check.Retry)
	a.Equal(3, check.Timeout)
}

func TestPublisherType(t *testing.T) {
	var cases = []struct {
		name   string
		reader io.Reader
	}{
		{"default", bytes.NewReader(defaultConfig)},
		{"notype", strings.NewReader("publishers: {memory: {}}")},
		{"caseinsensitive", strings.NewReader("publishers: {memory: {type: MEMORYPUBLISHER}}")},
	}

	for _, testcase := range cases {
		cfg, err := ReadConfig(testcase.reader)
		if err != nil {
			t.Fatalf("ReadConfig returned error for case %q: %s", testcase.name, err)
		}

		_, ok := cfg.Publishers["memory"]
		if !ok {
			t.Errorf("Expected 'memory' publisher not found (case %q)", testcase.name)
			continue
		}
		v, ok := cfg.Publishers["memory"]["type"]
		if !ok {
			t.Errorf("Expected 'type' element under 'memory' publisher (case %q)", testcase.name)
		}
		switch publisherType := v.(type) {
		case string:
			if publisherType != "memorypublisher" {
				t.Errorf("Expected 'type' to have value 'memorypublisher', not %q (case %q)", publisherType, testcase.name)
			}
		default:
			t.Errorf("Expected 'type' element to be a string, not %T (case %q)", publisherType, testcase.name)
		}
	}
}

func TestPublisherWrongType(t *testing.T) {
	_, err := ReadConfig(strings.NewReader("publishers: {memory: {type: []}}"))
	if errt, ok := err.(Error); !ok || errt.Code != ErrInvalidConfig {
		t.Error("Expected error with code ErrInvalidConfig")
	}

	_, err = ReadConfig(strings.NewReader("publishers: { invalid: {} }"))
	if errt, ok := err.(Error); !ok || errt.Code != ErrInvalidConfig {
		t.Error("Expected error with code ErrInvalidConfig")
	}
}

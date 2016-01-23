package nca

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
	"reflect"
	//	"strings"
	"testing"
)

func TestStartStop(t *testing.T) {
	r, err := NewRunnerFromFile("testdata/config.yml")
	require.Nil(t, err)
	r.Start()
	r.Stop()
}

func TestPublisherInit(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	runner, err := NewRunnerFromFile("testdata/config.yml")
	r.Nil(err)
	r.Equal(1, len(runner.publishers))

	expectType := reflect.TypeOf(&MemoryPublisher{})
	realType := reflect.TypeOf(runner.publishers["memory"])
	a.Equal(expectType, realType)
}

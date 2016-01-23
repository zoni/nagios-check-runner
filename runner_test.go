package nca

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
	"reflect"
	//	"strings"
	"testing"
	"time"
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

func TestRunAndPublishChecksCriticalCheck(t *testing.T) {
	runner, err := NewRunnerFromFile("testdata/criticalcheck.yml")
	require.Nil(t, err)
	publisher := runner.publishers["memory"]
	p := publisher.(*MemoryPublisher)

	runner.Start()
	defer runner.Stop()
	time.Sleep(500 * time.Millisecond)
	if count := p.EventCount(); count != 0 {
		t.Errorf("Expected eventcount to be 0, not %v", count)
	}
	time.Sleep(550 * time.Millisecond)
	if count := p.EventCount(); count != 1 {
		// Checks are run at retry interval after startup, so should have
		// triggered after 1 second
		t.Errorf("Expected eventcount to be 1, not %v", count)
	}
	time.Sleep(1 * time.Second)
	if count := p.EventCount(); count != 2 {
		t.Errorf("Expected eventcount to be 2, not %v", count)
	}
}

func TestRunAndPublishChecksOkCheck(t *testing.T) {
	runner, err := NewRunnerFromFile("testdata/okcheck.yml")
	require.Nil(t, err)
	publisher := runner.publishers["memory"]
	p := publisher.(*MemoryPublisher)

	runner.Start()
	defer runner.Stop()
	time.Sleep(1500 * time.Millisecond)
	if count := p.EventCount(); count != 0 {
		t.Errorf("Expected eventcount to be 0, not %v", count)
	}
	time.Sleep(550 * time.Millisecond)
	if count := p.EventCount(); count != 1 {
		// Checks are run at retry interval after startup, so should have
		// triggered after 1 second
		t.Errorf("Expected eventcount to be 1, not %v", count)
	}
	time.Sleep(1 * time.Second)
	if count := p.EventCount(); count != 2 {
		// Should now be triggering every second
		t.Errorf("Expected eventcount to be 2, not %v", count)
	}
}

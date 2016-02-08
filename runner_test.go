package ncr

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

func TestRunAndPublishChecks(t *testing.T) {
	runner, err := NewRunnerFromFile("testdata/runnertest.yml")
	require.Nil(t, err)
	publisher := runner.publishers["memory"]
	p := publisher.(*MemoryPublisher)

	runner.Start()
	defer runner.Stop()

	time.Sleep(500 * time.Millisecond)
	if count := p.EventCount("criticalcheck"); count != 0 {
		t.Errorf("Expected criticalcheck eventcount to be 0, not %v", count)
	}
	if count := p.EventCount("okcheck"); count != 0 {
		t.Errorf("Expected okcheck eventcount to be 0, not %v", count)
	}

	time.Sleep(550 * time.Millisecond)
	if count := p.EventCount("okcheck"); count != 0 {
		// Checks are run at retry interval after startup, so should not
		// trigger until 2 seconds
		t.Errorf("Expected okcheck eventcount to be 0, not %v", count)
	}
	if count := p.EventCount("criticalcheck"); count != 1 {
		// Checks are run at retry interval after startup, so should have
		// triggered after 1 second
		t.Errorf("Expected criticalcheck eventcount to be 1, not %v", count)
	}

	time.Sleep(1 * time.Second)
	if count := p.EventCount("okcheck"); count != 1 {
		// Okcheck should now have triggered once too, as retry = 2s
		t.Errorf("Expected okcheck eventcount to be 1, not %v", count)
	}
	if count := p.EventCount("criticalcheck"); count != 2 {
		// Criticalcheck has retry = 1s so should have triggered twice now
		t.Errorf("Expected criticalcheck eventcount to be 2, not %v", count)
	}

	time.Sleep(1 * time.Second)
	if count := p.EventCount("okcheck"); count != 2 {
		// Last result for okcheck was OK, so should now be
		// triggering at interval (1s) rather than retry interval (2s)
		t.Errorf("Expected okcheck eventcount to be 2, not %v", count)
	}
	if count := p.EventCount("criticalcheck"); count != 3 {
		t.Errorf("Expected criticalcheck eventcount to be 3, not %v", count)
	}
}

func assertCheckResult(t *testing.T, expected *CheckResult, actual *CheckResult) {
}

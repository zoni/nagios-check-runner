package ncr

import (
	log "gopkg.in/inconshreveable/log15.v2"
)

var Log = log.New()

const (
	StateOk       = 0
	StateWarning  = 1
	StateCritical = 2
	StateUnknown  = 3
)

const (
	ErrInvalidConfig = iota
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// Checker schedules and executes checks to be run.
type Checker interface {
	Start() (chan *CheckResult, error)
	Stop() error
	RegisterChecks(checks map[string]Check)
}

// Publisher publishes CheckResults.
type Publisher interface {
	// Start tells a Publisher to start so it can begin accepting check results.
	// It should be called after Configure.
	Start() error

	// Stop tells a Publisher to shut down.
	// No more check results may be published to it after calling Stop.
	Stop() error

	// Configure sets the configuration to be used by the Publisher.
	// It should be called before Start.
	Configure(cfg map[string]interface{}) error

	// Publish accepts a CheckResult to be published.
	// It should be safe for concurrent calling by multiple goroutines.
	Publish(*CheckResult) error
}

func init() {
	Log.SetHandler(log.DiscardHandler())
}

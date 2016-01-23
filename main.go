package nca

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
	Start() error
	Stop() error
	Configure(cfg map[string]interface{}) error
	Publish(*CheckResult) error
}

func init() {
	Log.SetHandler(log.DiscardHandler())
}

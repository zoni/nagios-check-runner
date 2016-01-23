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

func init() {
	Log.SetHandler(log.DiscardHandler())
}

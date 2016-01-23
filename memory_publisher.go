package nca

import (
	log "gopkg.in/inconshreveable/log15.v2"
)

type MemoryPublisher struct {
	log log.Logger
}

// Start starts this publisher.
func (p *MemoryPublisher) Start() error {
	p.log = Log.New("component", "memorypublisher")
	p.log.Info("Publisher ready")
	return nil
}

// Stop stops this publisher.
func (p *MemoryPublisher) Stop() error {
	p.log.Info("Publisher stopped")
	return nil
}

func (p *MemoryPublisher) Publish(result *CheckResult) error {
	return nil
}

// Configure sets the configuration to be used by the publisher.
func (p *MemoryPublisher) Configure(cfg map[string]interface{}) error {
	return nil
}

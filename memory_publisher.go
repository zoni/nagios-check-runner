package nca

import (
	log "gopkg.in/inconshreveable/log15.v2"
	"sync"
)

type MemoryPublisher struct {
	sync.Mutex
	log    log.Logger
	events []*CheckResult
}

// Start starts this publisher.
func (p *MemoryPublisher) Start() error {
	p.log = Log.New("component", "memorypublisher")
	p.log.Info("Publisher ready")
	p.events = make([]*CheckResult, 0)
	return nil
}

// Stop stops this publisher.
func (p *MemoryPublisher) Stop() error {
	p.log.Info("Publisher stopped")
	return nil
}

func (p *MemoryPublisher) Publish(result *CheckResult) error {
	p.Lock()
	defer p.Unlock()
	p.events = append(p.events, result)
	return nil
}

func (p *MemoryPublisher) EventCount() int {
	p.Lock()
	defer p.Unlock()
	return len(p.events)
}

func (p *MemoryPublisher) GetEvent(index int) *CheckResult {
	p.Lock()
	defer p.Unlock()
	return p.events[index]
}

// Configure sets the configuration to be used by the publisher.
func (p *MemoryPublisher) Configure(cfg map[string]interface{}) error {
	return nil
}

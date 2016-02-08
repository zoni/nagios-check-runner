package ncr

import (
	log "gopkg.in/inconshreveable/log15.v2"
	"sync"
)

// MemoryPublisher is a simple Publisher which stores all check results in memory.
// This publisher is used for testing and is not meant to be used normally.
type MemoryPublisher struct {
	sync.Mutex
	log    log.Logger
	events map[string][]*CheckResult
}

// Start starts this publisher.
func (p *MemoryPublisher) Start() error {
	p.log = Log.New("component", "memorypublisher")
	p.log.Info("Publisher ready")
	p.events = make(map[string][]*CheckResult)
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
	n := result.Name
	p.events[n] = append(p.events[n], result)
	return nil
}

func (p *MemoryPublisher) EventCount(checkname string) int {
	p.Lock()
	defer p.Unlock()
	return len(p.events[checkname])
}

func (p *MemoryPublisher) GetEvent(checkname string, index int) *CheckResult {
	p.Lock()
	defer p.Unlock()
	return p.events[checkname][index]
}

// Configure sets the configuration to be used by the publisher.
func (p *MemoryPublisher) Configure(cfg map[string]interface{}) error {
	return nil
}

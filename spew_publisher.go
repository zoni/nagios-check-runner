package nca

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	log "gopkg.in/inconshreveable/log15.v2"
)

type SpewPublisherConfig struct {
}

// SpewPublisher is a publisher which simply spews check results to stdout.
// This publisher is mainly useful for debugging
// and to use as a simple reference implementation for a Publisher.
type SpewPublisher struct {
	config SpewPublisherConfig
	log    log.Logger
}

// Start starts this publisher.
func (p *SpewPublisher) Start() error {
	p.log = Log.New("component", "spewpublisher")
	p.log.Info("Publisher ready")
	return nil
}

// Stop stops this publisher.
func (p *SpewPublisher) Stop() error {
	p.log.Info("Publisher stopped")
	return nil
}

// Publish prints the result of a check to stdout.
func (p *SpewPublisher) Publish(result *CheckResult) error {
	fmt.Println("\nReceived a check result:\n")
	spew.Dump(result)
	return nil
}

// Configure sets the configuration to be used by the publisher.
func (p *SpewPublisher) Configure(cfg map[string]interface{}) error {
	var result SpewPublisherConfig
	if err := mapstructure.Decode(cfg, &result); err != nil {
		return err
	}
	p.config = result
	return nil
}

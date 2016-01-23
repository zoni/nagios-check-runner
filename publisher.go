package nca

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	log "gopkg.in/inconshreveable/log15.v2"
)

type PublisherConfig struct {
	Foo string
}

// publisher is responsible for consuming check results from the checker
// and publishing them to a remote service.
type publisher struct {
	config PublisherConfig
	log    log.Logger
}

func (p *publisher) Start() error {
	p.log = Log.New("component", "publisher")
	p.log.Info("Publisher ready")
	return nil
}

func (p *publisher) Stop() error {
	p.log.Info("Publisher stopped")
	return nil
}

func (p *publisher) Publish(result *CheckResult) error {
	fmt.Println("Received a check result:")
	spew.Dump(result)
	return nil
}

func (p *publisher) SetConfig(cfg map[string]interface{}) error {
	var result PublisherConfig
	if err := mapstructure.Decode(cfg, &result); err != nil {
		return err
	}
	p.config = result
	return nil
}

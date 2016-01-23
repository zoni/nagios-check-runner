package nca

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "gopkg.in/inconshreveable/log15.v2"
)

// publisher is responsible for consuming check results from the checker
// and publishing them to a remote service.
type publisher struct {
	log log.Logger
}

//// run runs the publisher, consuming CheckResults from the publish channel.
//func (p *publisher) run() {
//for {
//select {
//case result := <-p.publish:
//fmt.Println("Received a check result:")
//spew.Dump(result)
//case <-p.done:
//fmt.Println("Publisher exiting")
//return
//}
//}
//}

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

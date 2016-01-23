package nca

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "gopkg.in/inconshreveable/log15.v2"
)

// publisher is responsible for consuming check results from the checker
// and publishing them to a remote service.
type publisher struct {
	publish chan *checkResult // The publish channel exposed by the Runner
	done    chan struct{}     // Used for signalling goroutines that we're shutting down
	log     log.Logger
}

// run runs the publisher, consuming checkResults from the publish channel.
func (p *publisher) run() {
	for {
		select {
		case result := <-p.publish:
			fmt.Println("Received a check result:")
			spew.Dump(result)
		case <-p.done:
			fmt.Println("Publisher exiting")
			break
		}
	}
}

func (p *publisher) Start() error {
	p.log = Log.New("component", "publisher")
	p.log.Info("Publisher starting")
	p.done = make(chan struct{})
	go p.run()
	return nil
}

func (p *publisher) Stop() error {
	p.log.Info("Publisher stopping")
	close(p.done)
	return nil
}

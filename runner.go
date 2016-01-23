package nca

import (
	//"github.com/davecgh/go-spew/spew"
	log "gopkg.in/inconshreveable/log15.v2"
	"os"
	"os/signal"
)

// Runner is the main component of the application. It runs the Checker and
// Publisher subcomponents and facilitates communication between them.
type Runner struct {
	config    Config
	publish   chan *checkResult
	checker   *checker
	publisher *publisher
	log       log.Logger
}

// NewRunner creates a new Runner with the given configuration.
func NewRunner(cfg Config) (*Runner, error) {
	r := &Runner{}
	if err := r.Init(cfg); err != nil {
		return nil, err
	}
	return r, nil
}

// NewRunnerFromFile creates a new Runner using the configuration loaded
// from the given file.
func NewRunnerFromFile(filename string) (*Runner, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg, err := LoadConfig(f)
	if err != nil {
		return nil, err
	}
	r, err := NewRunner(*cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Init initializes the Runner with the given configuration.
func (r *Runner) Init(cfg Config) error {
	r.config = cfg
	r.log = Log.New("component", "runner")
	r.publish = make(chan *checkResult)
	r.checker = &checker{publish: r.publish, checks: r.config.Checks}
	r.publisher = &publisher{publish: r.publish}
	return nil
}

// Start starts the Runner and begins running checks, publishing them as they
// complete.
func (r *Runner) Start() error {
	r.log.Info("Runner starting")
	r.checker.Start()
	r.publisher.Start()
	return nil
}

// Stop stops and shuts down the Runner.
func (r *Runner) Stop() error {
	r.log.Info("Runner stopping")
	r.checker.Stop()
	r.publisher.Stop()
	return nil
}

// Run starts the runner and blocks until an interrupt signal is received.
func (r *Runner) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	r.Start()
	for sig := range c {
		if sig == os.Interrupt {
			r.log.Info("Received interrupt, shutting down")
			r.Stop()
			break
		}
	}

	return nil
}

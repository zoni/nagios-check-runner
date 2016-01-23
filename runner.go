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
	config      Config
	checker     *checker
	publishers  map[string]Publisher
	publishChan chan *CheckResult // The cchannel check results are received from
	log         log.Logger
	done        chan struct{} // Used for signalling goroutines that we're shutting down
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
	cfg, err := ReadConfig(f)
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
	r.checker = &checker{}
	r.checker.RegisterChecks(r.config.Checks)

	r.publishers = make(map[string]Publisher)
	for label, config := range cfg.Publishers {
		p := &SpewPublisher{}
		if err := p.Configure(config); err != nil {
			r.log.Error("Invalid publisher configuration", "publisher", label, "error", err)
			return err
		}
		r.publishers[label] = p
	}

	return nil
}

// Start starts the Runner and begins running checks, publishing them as they
// complete.
func (r *Runner) Start() error {
	r.log.Info("Runner starting")
	r.done = make(chan struct{})

	r.publishChan, _ = r.checker.Start()
	for _, publisher := range r.publishers {
		publisher.Start()
	}
	go r.process()
	return nil
}

// process reads results produced by the checker and distributes them
// to the publishers.
func (r *Runner) process() {
	for {
		select {
		case result := <-r.publishChan:
			for _, publisher := range r.publishers {
				publisher.Publish(result)
			}
		case <-r.done:
			return
		}
	}
}

// Stop stops and shuts down the Runner.
func (r *Runner) Stop() error {
	r.log.Info("Runner stopping")
	r.checker.Stop()
	for _, publisher := range r.publishers {
		publisher.Stop()
	}
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

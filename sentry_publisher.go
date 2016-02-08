package ncr

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/mitchellh/mapstructure"
	log "gopkg.in/inconshreveable/log15.v2"
	"os"
	"strconv"
)

// SentryPublisherConfig describes the configuration used by the SentryPublisher.
type SentryPublisherConfig struct {
	Dsn      string // The project DSN to use
	Hostname string // The hostname to report to sentry
}

// SentryPublisher is a Publisher which publishes non-OK check results to Sentry
// (https://getsentry.com/)
type SentryPublisher struct {
	config SentryPublisherConfig
	log    log.Logger
	client *raven.Client
}

// Start starts this publisher.
func (p *SentryPublisher) Start() error {
	p.log = Log.New("component", "sentrypublisher")
	client, err := raven.New(p.config.Dsn)
	if err != nil {
		p.log.Error("Sentry client failed to initialize", "error", err)
		return err
	}
	p.client = client
	p.log.Info("Publisher ready")
	return nil
}

// Stop stops this publisher.
func (p *SentryPublisher) Stop() error {
	p.client.Wait()
	p.log.Info("Publisher stopped")
	return nil
}

// Publish publishes the CheckResult to Sentry if it has a state other than StateOk.
func (p *SentryPublisher) Publish(result *CheckResult) error {
	if result.Returncode == StateOk {
		p.log.Debug("Not publishing to Sentry", "reason", "State is OK")
		return nil
	}

	var state string
	switch result.Returncode {
	case StateWarning:
		state = "warning"
	case StateCritical:
		state = "critical"
	default:
		state = "unknown"
	}

	p.client.CaptureMessageAndWait(
		fmt.Sprintf("%q is %s on %s", result.Name, state, p.config.Hostname),
		map[string]string{
			"name":        result.Name,
			"returncode":  strconv.Itoa(result.Returncode),
			"output":      string(result.Output),
			"server_name": p.config.Hostname,
		},
	)
	return nil
}

// Configure sets the configuration to be used by the publisher.
func (p *SentryPublisher) Configure(cfg map[string]interface{}) error {
	var result SentryPublisherConfig
	if err := mapstructure.Decode(cfg, &result); err != nil {
		return err
	}
	if result.Dsn == "" {
		return errors.New("dsn must be specified")
	}
	if result.Hostname == "" {
		h, _ := os.Hostname()
		result.Hostname = h
	}
	p.config = result
	return nil
}

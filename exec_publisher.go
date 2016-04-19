package ncr

import (
	"bytes"
	"errors"
	"github.com/kballard/go-shellquote"
	"github.com/mitchellh/mapstructure"
	log "gopkg.in/inconshreveable/log15.v2"
	"os/exec"
	"text/template"
)

// ExecPublisherConfig describes the configuration used by the ExecPublisher.
type ExecPublisherConfig struct {
	Cmd   string // The command to execute on Publish
	Stdin string // Template to use to format stdin to Cmd
}

// ExecPublisher is a Publisher which executes an external program whenever
// a check result is received.
type ExecPublisher struct {
	config        ExecPublisherConfig
	log           log.Logger
	cmd           []string // The command to execute on publishing
	stdinTemplate *template.Template
}

// Start starts this publisher.
func (p *ExecPublisher) Start() error {
	p.log = Log.New("component", "execpublisher")
	p.log.Info("Publisher ready")
	return nil
}

// Stop stops this publisher.
func (p *ExecPublisher) Stop() error {
	p.log.Info("Publisher stopped")
	return nil
}

// Publish prints the result of a check to stdout.
func (p *ExecPublisher) Publish(result *CheckResult) error {
	var cmd *exec.Cmd
	if len(p.cmd) > 1 {
		cmd = exec.Command(p.cmd[0], p.cmd[1:]...)
	} else {
		cmd = exec.Command(p.cmd[0])
	}
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		p.log.Error("Exec failed", "error", err)
		return err
	}
	if err = p.stdinTemplate.Execute(stdin, result); err != nil {
		p.log.Error("Failed to write template data to stdin", "error", err)
		cmd.Wait() // Avoid zombies
		return err
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		p.log.Error("Exec failed", "error", err, "output", b.String())
		return err
	}
	return nil
}

// Configure sets the configuration to be used by the publisher.
func (p *ExecPublisher) Configure(cfg map[string]interface{}) error {
	var result ExecPublisherConfig
	if err := mapstructure.Decode(cfg, &result); err != nil {
		return err
	}
	if result.Cmd == "" {
		return errors.New("cmd must be specified")
	}
	args, err := shellquote.Split(result.Cmd)
	if err != nil {
		return err
	}
	p.cmd = args

	tmpl, err := template.New("stdin").Parse(result.Stdin)
	if err != nil {
		return nil
	}
	p.stdinTemplate = tmpl
	p.config = result
	return nil
}

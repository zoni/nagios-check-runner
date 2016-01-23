package main

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gophergala2016/nagios-check-runner"
	log "gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/term"
	"os"
)

func init() {
	var handler log.Handler
	if term.IsTty(os.Stdout.Fd()) {
		handler = log.StreamHandler(os.Stdout, log.TerminalFormat())
	} else {
		handler = log.StreamHandler(os.Stdout, log.JsonFormat())
	}
	log.Root().SetHandler(handler)
	nca.Log.SetHandler(handler)
}

func main() {
	r, err := nca.NewRunnerFromFile("config.yml")
	if err != nil {
		log.Crit("Startup failed", "error", err)
		os.Exit(1)
	}

	r.Run()
}

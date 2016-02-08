package main

import (
	"fmt"
	"github.com/zoni/nagios-check-runner"
	"gopkg.in/alecthomas/kingpin.v2"
	log "gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/term"
	"os"
	"runtime"
)

var (
	version = "<unknown version/release>"

	app     = kingpin.New("ncr", "Nagios Check Runner")
	cfgfile = app.Flag("config", "Configuration file to load settings from.").Short('c').Default("/etc/ncr/ncr.yml").String()
)

func init() {
	var handler log.Handler
	if term.IsTty(os.Stdout.Fd()) {
		handler = log.StreamHandler(os.Stdout, log.TerminalFormat())
	} else {
		handler = log.StreamHandler(os.Stdout, log.JsonFormat())
	}
	log.Root().SetHandler(handler)
	ncr.Log.SetHandler(handler)
}

func main() {
	app.Version(fmt.Sprintf("NCR version %s (runtime version %s)", version, runtime.Version()))
	kingpin.MustParse(app.Parse(os.Args[1:]))

	r, err := ncr.NewRunnerFromFile(*cfgfile)
	if err != nil {
		log.Crit("Startup failed", "error", err)
		os.Exit(1)
	}

	r.Run()
}

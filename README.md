Nagios check runner
===================


*Run nagios-compatible plugins locally and publish their results anywhere you can think of.*

Nagios check runner was originally conceived as a replacement for
[NRPE](https://exchange.nagios.org/directory/Addons/Monitoring-Agents/NRPE--2D-Nagios-Remote-Plugin-Executor/details),
designed to run checks locally and publish their results as passive checks using
[NSCA-ng](https://www.nsca-ng.org/).

It has grown to support multiple output plugins, potentially making it useful
even without a traditional Nagios/Icinga setup. Small one-person shop with just
one or two servers to monitor? Receive reports by email when services go down or
transmit them to [Sentry](https://getsentry.com).

Using a PAAS monitoring solution but want to tap into the thousands of Nagios
plugins already available? Publish results to a webhook as HTTP events.

### Submitting passive check results to Icinga using NSCA-ng

![Screencast of submitting passive check results using nsca-ng](screencast/nsca.gif)

*Note: If you look closely you'll see a delay in between the check result being submitted
and the service state updating in Icinga. This is not a bug in the runner but a delay in
NSCA-ng submitting the event to Icinga and Icinga's event loop picking up the change.*

### Submitting check results to Sentry

![Screencast of nagios plugins publishing to Sentry](screencast/sentry.gif)


Configuration
-------------

Configuration is split up in two parts, *checks* and *publishers*. For example:

``` yaml
publishers:
  nsca_ng:
    type: ExecPublisher
    cmd: /usr/sbin/send_nsca
    stdin: |
        "myhost	{{ .Name }}	{{ .Returncode }}	{{ .Output | printf "%s" }}

checks:
  Local SMTP:
    command: /usr/lib/nagios/plugins/check_smtp -H localhost
    interval: 30
    retry: 10
    timeout: 3

  System load:
    command: /usr/lib/nagios/plugins/check_load -r -w 3,2,2 -c 5,3,3
    interval: 60
    retry: 60
    timeout: 10
```

___Checks___ define the checks to be run. The title is the service name and the
options it takes are as follows:

* *command*: The command to be run.
* *interval*: The interval to run this plugin when it is in OK state.
* *retry*: The interval to run this plugin when it returns a state other
than OK.
* *timeout*: Kills the process if it does not finish within this time limit.

___Publishers___ defines the publishers to which check results should be published.
The available publishers and their options are detailed below. All entries under
*publishers* take a descriptive name of your own choosing and must have a `type`
element listing the publisher to be used from the list below.

### SpewPublisher

This is a debug publisher, useful for local development. It prints check results
to stdout.

### ExecPublisher

Executes an external command for every received result. This allows for complete
freedom in what to do with check results and can be used to publish results as
passive checks using [nsca-ng](https://www.nsca-ng.org/), as email by feeding it
to `mail`/`sendmail` or some other custom script of your choosing.

Options:

* *cmd*: The path to the script/binary to be run.
* *stdin*: A text template describing what data is to be fed as input over
  stdin to the specified cmd. It uses Go
  [text/template](https://golang.org/pkg/text/template/) templating.

  Available template values are `.Name`, `.Returncode`, and `.Output`.

### SentryPublisher

Publishes results as messages to the exception tracking system
[Sentry](https://getsentry.com/).

Options:

* *dsn*: The DSN of the project to submit to.
* *hostname (optional)*: The hostname to report to sentry.


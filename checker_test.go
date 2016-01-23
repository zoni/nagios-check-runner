package nca

import (
	"bytes"
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
)

func TestRunCheck(t *testing.T) {
	a := assert.New(t)

	check := Check{
		Name:     "Dummy",
		Command:  "",
		Args:     []string{"testdata/check_dummy", "0", "OK"},
		Interval: 1,
		Retry:    1,
		Timeout:  5,
	}
	result := runCheck(check)
	a.Equal(StateOk, result.Returncode)

	check.Args = []string{"testdata/check_dummy", "2", "CRITICAL!"}
	result = runCheck(check)
	a.Equal(StateCritical, result.Returncode)
	a.Equal(
		bytes.TrimSpace([]byte("CRITICAL!")),
		bytes.TrimSpace(result.Output),
	)

	check.Args = []string{"/no/such/file/or/directory"}
	result = runCheck(check)
	a.Equal(StateCritical, result.Returncode)

	// Note: Leave this test last or set timeout back
	check.Args = []string{"/bin/sleep", "2"}
	check.Timeout = 1
	result = runCheck(check)
	a.Equal(StateCritical, result.Returncode)
	a.Equal(
		bytes.TrimSpace([]byte("Process /bin/sleep was killed after 1 second timeout")),
		bytes.TrimSpace(result.Output),
	)
}

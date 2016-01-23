package nca

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
	"io/ioutil"
)

func TestConfigParsing(t *testing.T) {
	assert := assert.New(t)
}

func writeTempConfig(t *testing.T, config []byte) {
	f, err := ioutil.TempFile(dir, prefix string) (f *os.File, err error)
}

package nca

import (
	//"bytes"
	//"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
	//"io/ioutil"
	//"os"
	//"strings"
	"testing"
)

//var defaultConfig []byte

//func init() {
//f, err := os.Open("testdata/config.yml")
//if err != nil {
//panic(err)
//}
//defer f.Close()

//defaultConfig, err = ioutil.ReadAll(f)
//if err != nil {
//panic(err)
//}
//}

func TestStartStop(t *testing.T) {
	//a := assert.New(t)

	r, err := NewRunnerFromFile("testdata/config.yml")
	require.Nil(t, err)
	r.Start()
	r.Stop()
}

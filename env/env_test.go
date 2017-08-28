package env

import (
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

type Env struct {
	Addrs          []string      `envconfig:"ADDRS" default:"localhost:9042"`
	TimeoutConnect time.Duration `envconfig:"TIMEOUT_CONNECT" default:"30s"`
	Timeout        time.Duration `envconfig:"TIMEOUT" default:"30s"`
	NbConns        int           `envconfig:"NB_CONNS" default:"8"`
	HostFilter     string        `envconfig:"HOST_FILTER" default:""`
	TimeoutLimit   int64         `envconfig:"TIMEOUT_LIMIT" default:"10"`
}

func TestEnv_String(t *testing.T) {
	assert.Equal(t, "", String(42))

	e := &Env{}
	assert.NoError(t, envconfig.Process("", e))
	expected := "ADDRS = [localhost:9042]\nTIMEOUT_CONNECT = 30s\nTIMEOUT = 30s\nNB_CONNS = 8\nHOST_FILTER = \nTIMEOUT_LIMIT = 10"
	assert.Equal(t, expected, String(e))
}

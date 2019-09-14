package modbus

import (
	"testing"

	cfg "github.com/mythay/lark/config"
	"github.com/stretchr/testify/assert"
)

func Test_agent_step(t *testing.T) {
	assert := assert.New(t)
	c, err := cfg.LoadCfgModbus("testdata/modbus_sample.yaml")

	assert.NoError(err)

	agnt, err := newAgent(*c)

	assert.NoError(err)

	agnt.step()
}

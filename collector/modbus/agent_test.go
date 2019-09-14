package modbus

import (
	"context"
	"os/exec"
	"testing"
	"time"

	cfg "github.com/mythay/lark/config"
	"github.com/stretchr/testify/assert"
)

func Test_agent_step(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	err := exec.CommandContext(ctx, "testdata/diagslave", "-m", "tcp", "-p", "5020").Start()
	assert.NoError(err)

	c, err := cfg.LoadCfgModbus("testdata/modbus_sample.yaml")

	assert.NoError(err)

	agnt, err := newAgent(*c)

	assert.NoError(err)

	agnt.step()
}

func Test_agent_run(t *testing.T) {
	assert := assert.New(t)
	// ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	// defer cancel()
	// err := exec.CommandContext(ctx, "testdata/diagslave", "-m", "tcp", "-p", "5020").Start()
	// assert.NoError(err)

	c, err := cfg.LoadCfgModbus("testdata/modbus_sample.yaml")

	assert.NoError(err)

	agnt, err := newAgent(*c)

	assert.NoError(err)
	t.Logf("ready to start")

	agnt.Start()

	time.Sleep(5 * time.Second)
	agnt.Stop()
	time.Sleep(1 * time.Second)
}

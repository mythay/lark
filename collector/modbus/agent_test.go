package modbus

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func metricPrint(measurement string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {
	fmt.Printf("[%s] %v-%v\n", measurement, tags, fields)
}
func Test_agent_step(t *testing.T) {
	assert := assert.New(t)
	var err error
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	err = exec.CommandContext(ctx, "testdata/diagslave", "-m", "tcp", "-p", "5020").Start()
	assert.NoError(err)

	agent, err := NewAgent()

	assert.NoError(err)
	err = agent.LoadConfig("testdata/modbus_sample.yaml")

	assert.NoError(err)
	agent.Gather(metricPrint)
}

func Test_agent_run(t *testing.T) {
	assert := assert.New(t)
	// ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	// defer cancel()
	// err := exec.CommandContext(ctx, "testdata/diagslave", "-m", "tcp", "-p", "5020").Start()
	// assert.NoError(err)

	agent, err := NewAgent()

	assert.NoError(err)

	err = agent.LoadConfig("testdata/modbus_sample.yaml")
	assert.NoError(err)

	agent.Start(func(measurement string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {})

}

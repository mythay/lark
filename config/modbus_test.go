package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTomlConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := LoadCfgModbus("testdata/modbus_sample.toml")
	assert.Empty(err)

	assert.Equal("modbus", config.Type)
	assert.Equal(1, len(config.Products))

	assert.Equal(3, len(config.Products["em3250"].Registers))
	assert.Equal("voltage", config.Products["em3250"].Registers["input-1"].Tag)
	assert.Equal(2, len(config.Products["em3250"].Ranges))

}

func TestRegTypeAndSize(t *testing.T) {

	assert := assert.New(t)

	reg := CfgRegister{Type: "int16"}
	assert.Equal(reflect.TypeOf(int16(0)), reg.GoType())
	assert.EqualValues(1, reg.Count())
}

package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYamlConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := LoadCfgModbus("testdata/modbus_sample.yaml")
	assert.Empty(err)

	assert.Equal("modbus", config.Type)
	assert.Equal(1, len(config.Profile))

	profile := config.Profile[0]

	assert.Equal(2, len(profile.Register))
	assert.Equal(3, len(profile.Range))

	reg := profile.Register[0]
	assert.Equal("voltage", reg.Catalog)

}

func TestRegTypeAndSize(t *testing.T) {

	assert := assert.New(t)

	reg := CfgRegister{Type: "int16"}
	assert.Equal(reflect.TypeOf(int16(0)), reg.GoType())
	assert.EqualValues(1, reg.Count())
}
func TestArrayRange(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	for i, a := range arr {
		fmt.Printf("%v, %v, %v\n", a, &a, &arr[i])
	}

}

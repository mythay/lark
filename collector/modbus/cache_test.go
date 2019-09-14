package modbus

import (
	"sort"
	"testing"

	"github.com/mythay/lark/collector/modbus/mocks"
	cfg "github.com/mythay/lark/config"

	"github.com/stretchr/testify/assert"
)

//go:generate mockery -name=MbReader

func Test_Range_Sort(t *testing.T) {
	assert := assert.New(t)
	slice := rangeSlice{
		{org: cfg.CfgRange{Start: 10, End: 20}},
		{org: cfg.CfgRange{Start: 1, End: 10}},
		{org: cfg.CfgRange{Start: 20, End: 30}},
	}
	sort.Sort(slice)
	assert.Equal(uint16(1), slice[0].org.Start)
	assert.Equal(uint16(10), slice[1].org.Start)
	assert.Equal(uint16(20), slice[2].org.Start)

}

var rags = []*cfg.CfgRange{
	&cfg.CfgRange{Start: 10, End: 20},
	&cfg.CfgRange{Start: 1, End: 10},
	&cfg.CfgRange{Start: 20, End: 30, Fixed: true},
}
var regs = []*cfg.CfgRegister{
	&cfg.CfgRegister{Start: 1, Quantity: 1},
	&cfg.CfgRegister{Start: 3, Quantity: 4},
	&cfg.CfgRegister{Start: 15, Quantity: 2},
	&cfg.CfgRegister{Start: 10, Quantity: 2},
	&cfg.CfgRegister{Start: 20, Quantity: 2},
}

func Test_mbCache_check_range(t *testing.T) {
	assert := assert.New(t)
	mk := &mocks.MbReader{}
	cache, err := newCache(rags, regs, mk)
	assert.NoError(err)

	assert.NotNil(cache.locateRange(cfg.CfgRange{Start: 1, End: 2}))
	assert.NotNil(cache.locateRange(cfg.CfgRange{Start: 1, End: 10}))
	assert.NotNil(cache.locateRange(cfg.CfgRange{Start: 20, End: 20}))
	assert.NotNil(cache.locateRange(cfg.CfgRange{Start: 20, End: 30}))

	assert.Nil(cache.locateRange(cfg.CfgRange{Start: 0, End: 1}))
	assert.Nil(cache.locateRange(cfg.CfgRange{Start: 1, End: 11}))
	assert.Nil(cache.locateRange(cfg.CfgRange{Start: 10, End: 21}))
	assert.Nil(cache.locateRange(cfg.CfgRange{Start: 10, End: 21}))
	assert.Nil(cache.locateRange(cfg.CfgRange{Start: 41, End: 41}))

}

func Test_mbCache_read_with_cache(t *testing.T) {
	assert := assert.New(t)
	mk := &mocks.MbReader{}
	cache, err := newCache(rags, regs, mk)
	assert.NoError(err)
	assert.NotNil(cache)
	mk.On("ReadHolding", uint16(1), uint16(6)).Return([]byte{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6}, nil).Once()
	data, err := cache.Read(1, 1)
	assert.Equal(uint16(1), data.quantity())
	assert.Equal([]uint8{0x1, 0x1}, data.value)

	data, err = cache.Read(2, 1)
	assert.Equal(uint16(1), data.quantity())
	assert.Equal([]uint8{0x2, 0x2}, data.value)

	mk.AssertNumberOfCalls(t, "ReadHolding", 1)

}

func Test_mbCache_read_with_cache_and_fixed_range(t *testing.T) {
	assert := assert.New(t)
	mk := &mocks.MbReader{}
	cache, err := newCache(rags, regs, mk)
	assert.NoError(err)
	assert.NotNil(cache)
	// mk.AssertCalled(t, "ReadHolding", uint16(1), uint16(6))
	mk.On("ReadHolding", uint16(20), uint16(11)).Return([]byte{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11}, nil).Once()
	data, err := cache.Read(20, 2)
	assert.Equal(uint16(2), data.quantity())
	assert.Equal([]uint8{0x1, 0x1, 2, 2}, data.value)

	data, err = cache.Read(21, 2)
	assert.Equal(uint16(2), data.quantity())
	assert.Equal([]uint8{0x2, 0x2, 3, 3}, data.value)

	mk.AssertNumberOfCalls(t, "ReadHolding", 1)

}

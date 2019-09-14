package modbus

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalEndian_Decode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint16(0x0102), NormalEndian.Uint16([]byte{1, 2}))
	assert.Equal(uint32(50594050), NormalEndian.Uint32([]byte{1, 2, 3, 4}))
	assert.Equal(uint64(0x0708050603040102), NormalEndian.Uint64([]byte{1, 2, 3, 4, 5, 6, 7, 8}))

}

func TestNormalEndian_Encode(t *testing.T) {
	assert := assert.New(t)

	var data [8]byte

	NormalEndian.PutUint16(data[:2], 0x0102)
	assert.EqualValues(data[:2], []byte{1, 2})

	NormalEndian.PutUint32(data[:4], 50594050)
	assert.EqualValues(data[:4], []byte{1, 2, 3, 4})

	NormalEndian.PutUint64(data[:], 0x0708050603040102)
	assert.EqualValues(data[:], []byte{1, 2, 3, 4, 5, 6, 7, 8})
	// assert.Equal(NormalEndian.PutUint16(0x0102), NormalEndian.Uint16([]byte{1, 2}))
	// assert.Equal(uint32(50594050), NormalEndian.Uint32([]byte{1, 2, 3, 4}))
	// assert.Equal(uint64(0x0708050603040102), NormalEndian.Uint64([]byte{1, 2, 3, 4, 5, 6, 7, 8}))

}

func TestNormalEndian_signed(t *testing.T) {
	assert := assert.New(t)

	var signed16 int16
	binary.Read(bytes.NewReader([]byte{0xff, 2}), NormalEndian, &signed16)
	assert.Equal(int16(-254), signed16)

	var signed32 int32
	binary.Read(bytes.NewReader([]byte{1, 2, 0xff, 4}), NormalEndian, &signed32)
	assert.Equal(int32(-16514814), signed32)

}

func TestNormalEndian_float(t *testing.T) {
	assert := assert.New(t)

	var fnum float32
	binary.Read(bytes.NewReader([]byte{0xff, 2, 0xff, 4}), NormalEndian, &fnum)
	assert.InDelta(float32(-1.76782e+038), fnum, .001e+38)

	var dnum float64
	binary.Read(bytes.NewReader([]byte{1, 2, 3, 4, 5, 6, 0xff, 8}), NormalEndian, &dnum)
	assert.InDelta(float64(-8.23591448547104e+303), dnum, .001e+303)

}

///

func TestInverseEndian_Decode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint16(0x0102), InverseEndian.Uint16([]byte{1, 2}))
	assert.Equal(uint32(16909060), InverseEndian.Uint32([]byte{1, 2, 3, 4}))
	assert.Equal(uint64(0x0102030405060708), InverseEndian.Uint64([]byte{1, 2, 3, 4, 5, 6, 7, 8}))

}

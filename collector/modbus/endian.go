package modbus

var NormalEndian normalEndian

type normalEndian struct{}

func (normalEndian) Uint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

func (normalEndian) PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[1] = byte(v)
	b[0] = byte(v >> 8)
}

func (normalEndian) Uint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[1]) | uint32(b[0])<<8 | uint32(b[3])<<16 | uint32(b[2])<<24
}

func (normalEndian) PutUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[1] = byte(v)
	b[0] = byte(v >> 8)
	b[3] = byte(v >> 16)
	b[2] = byte(v >> 24)
}

func (normalEndian) Uint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[1]) | uint64(b[0])<<8 | uint64(b[3])<<16 | uint64(b[2])<<24 |
		uint64(b[5])<<32 | uint64(b[4])<<40 | uint64(b[7])<<48 | uint64(b[6])<<56
}

func (normalEndian) PutUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[1] = byte(v)
	b[0] = byte(v >> 8)
	b[3] = byte(v >> 16)
	b[2] = byte(v >> 24)
	b[5] = byte(v >> 32)
	b[4] = byte(v >> 40)
	b[7] = byte(v >> 48)
	b[6] = byte(v >> 56)
}

func (normalEndian) String() string { return "NormalEndian" }

func (normalEndian) GoString() string { return "modbus.NormalEndian" }

var InverseEndian inverseEndian

type inverseEndian struct{}

func (inverseEndian) Uint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

func (inverseEndian) PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[1] = byte(v)
	b[0] = byte(v >> 8)
}

func (inverseEndian) Uint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func (inverseEndian) PutUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[3] = byte(v)
	b[2] = byte(v >> 8)
	b[1] = byte(v >> 16)
	b[0] = byte(v >> 24)
}

func (inverseEndian) Uint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func (inverseEndian) PutUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[7] = byte(v)
	b[6] = byte(v >> 8)
	b[5] = byte(v >> 16)
	b[4] = byte(v >> 24)
	b[3] = byte(v >> 32)
	b[2] = byte(v >> 40)
	b[1] = byte(v >> 48)
	b[0] = byte(v >> 56)
}

func (inverseEndian) String() string { return "InverseEndian" }

func (inverseEndian) GoString() string { return "modbus.InverseEndian" }

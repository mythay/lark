package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	cfg "github.com/mythay/lark/config"
)

type regValue struct {
	*cfg.CfgRegister
	Value interface{}
	Ts    time.Time
	name  string
}

func (reg *regValue) Endian() binary.ByteOrder {
	if reg.Inverse {
		return InverseEndian
	} else {
		return NormalEndian
	}
}

// performance is 10x slower than no reflect function, left it here just for reference
func (reg *regValue) ReflectParse(data []byte) (interface{}, error) {
	var err error
	val := reflect.New(reg.GoType())

	buf := bytes.NewReader(data)
	err = binary.Read(buf, reg.Endian(), val.Interface())
	return val.Elem().Interface(), err
}

func (reg *regValue) Parse(data []byte) (interface{}, error) {

	switch reg.GoKind() {
	case reflect.Uint16:
		return reg.Endian().Uint16(data), nil
	case reflect.Int16:
		return int16(reg.Endian().Uint16(data)), nil
	case reflect.Int32:
		return int32(reg.Endian().Uint32(data)), nil
	case reflect.Float32:
		return float32(reg.Endian().Uint32(data)), nil
	case reflect.Float64:
		return float64(reg.Endian().Uint64(data)), nil
	}
	return nil, fmt.Errorf("invalid type")

}

// parse response data and store it and time stamp to its own storage

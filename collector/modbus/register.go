package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	cfg "github.com/mythay/lark/config"
)

type mbReg struct {
	*cfg.CfgRegister
	Value interface{}
	Ts    time.Time
}

func (reg *mbReg) Endian() binary.ByteOrder {
	if reg.Inverse {
		return InverseEndian
	} else {
		return NormalEndian
	}
}

func (reg *mbReg) preCondition() bool {
	return true
}

func (reg *mbReg) postCondition() bool {
	return true
}

func (reg *mbReg) step(reader *mbCache) bool {
	if !reg.preCondition() {
		return false
	}
	data, err := reader.Read(reg.Start, reg.Quantity)
	if err != nil {
		return false
	}
	reg.Value, err = reg.parse(data.value)
	if err != nil {
		return false
	}
	if !reg.postCondition() {
		return false
	}
	return true
}

// performance is 10x slower than no reflect function, left it here just for reference
func (reg *mbReg) reflectParse(data []byte) (interface{}, error) {
	var err error
	val := reflect.New(reg.GoType())

	buf := bytes.NewReader(data)
	err = binary.Read(buf, reg.Endian(), val.Interface())
	return val.Elem().Interface(), err
}

func (reg *mbReg) parse(data []byte) (interface{}, error) {

	switch reg.GoKind() {
	case reflect.Uint16:
		return reg.Endian().Uint16(data), nil
	case reflect.Int16:
		return int16(reg.Endian().Uint16(data)), nil
	case reflect.Int32:
		return int32(reg.Endian().Uint32(data)), nil
	case reflect.Uint32:
		return uint32(reg.Endian().Uint32(data)), nil
	case reflect.Float32:
		return float32(reg.Endian().Uint32(data)), nil
	case reflect.Float64:
		return float64(reg.Endian().Uint64(data)), nil
	}
	return nil, fmt.Errorf("invalid type")

}

// parse response data and store it and time stamp to its own storage

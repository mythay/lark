package mbclient

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/mythay/spider"
)

type regValue struct {
	spider.CfgRegister
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
func (reg *regValue) ParseResponse(resp map[uint16]tsRegValue) {
	count := reg.Count()
	var found = true
	var ts time.Time
	rawdata := make([]byte, 8)
	for i := reg.Base; i < reg.Base+count; i++ {
		if d, ok := resp[i]; ok {
			rawdata = append(rawdata, d.value[0], d.value[1])
			ts = d.ts
		} else {
			found = false
			break
		}
	}
	if found {
		val, err := reg.Parse(rawdata)
		if err == nil {
			reg.Value = val
			reg.Ts = ts
			// fmt.Printf("%s=%v\n", reg.name, val)
		}
	}
}

package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"time"

	cfg "github.com/mythay/lark/config"
)

func parseTagString(query string) map[string]string {
	m := map[string]string{}
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "& "); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		if len(key) == 0 || len(value) == 0 {
			continue
		}

		m[key] = value
	}
	return m
}

type mbReg struct {
	*cfg.CfgRegister
	Value interface{}
	Ts    time.Time

	tags map[string]string
}

func NewReg(cfg *cfg.CfgRegister) (*mbReg, error) {
	reg := &mbReg{CfgRegister: cfg, tags: parseTagString(cfg.Tag)}
	reg.tags["catalog"] = cfg.Catalog
	return reg, nil
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

func (reg *mbReg) collect(reader *mbCache) bool {
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

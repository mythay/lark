package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/influxdata/toml"
)

type Verifier interface {
	Verify() (bool, error)
}

type CfgModbus struct {
	Type     string
	Products map[string]CfgProduct
	Hosts    []CfgHost
}
type CfgProduct struct {
	Registers map[string]CfgRegister
	Ranges    []CfgRange
	// Inverse  bool
}
type CfgRegister struct {
	Base    uint16
	Type    string
	Inverse bool
	Cmd     string
	Tag     string
	//
	cCount uint16
	cType  reflect.Type
	cKind  reflect.Kind
}

func (reg *CfgRegister) Count() uint16 {
	if reg.cCount == 0 {
		reg.cCount = uint16(reg.GoType().Size()) / 2
	}
	return reg.cCount
}

func (reg *CfgRegister) GoType() reflect.Type {
	if reg.cType == nil {
		switch reg.Type {
		case "int16":
			reg.cType = reflect.TypeOf(int16(0))
		case "uint16":
			reg.cType = reflect.TypeOf(uint16(0))
		case "int32":
			reg.cType = reflect.TypeOf(int32(0))
		case "uint32":
			reg.cType = reflect.TypeOf(uint32(0))
		case "float32":
			reg.cType = reflect.TypeOf(float32(0))
		case "float64":
			reg.cType = reflect.TypeOf(float64(0))
		default:
			return reflect.TypeOf(nil)
		}
	}
	return reg.cType
}

func (reg *CfgRegister) GoKind() reflect.Kind {
	if reg.cKind == reflect.Invalid {
		switch reg.Type {
		case "int16":
			reg.cKind = reflect.Int16
		case "uint16":
			reg.cKind = reflect.Uint16
		case "int32":
			reg.cKind = reflect.Int32
		case "uint32":
			reg.cKind = reflect.Uint32
		case "float32":
			reg.cKind = reflect.Float32
		case "float64":
			reg.cKind = reflect.Float64
		default:
			return reflect.Invalid
		}
	}
	return reg.cKind
}

func (reg *CfgRegister) Verify() (bool, error) {
	if reg.Count() != 0 {
		return true, nil
	} else {
		return false, fmt.Errorf(" invalid type'%s'", reg.Type)
	}

}
func (reg *CfgRegister) Last() uint16 {
	return reg.Base + reg.Count()
}

type CfgRange struct {
	Start uint16
	End   uint16
}

func (rag *CfgRange) Count() uint16 {
	return rag.End - rag.Start + 1
}

type CfgAddrTcp struct {
	IpAddr string
	Port   uint16
}

type CfgAddrRtu struct {
	SerialPort string
	BaudRate   int
	DataBits   int    // Data bits: 5, 6, 7 or 8 (default 8)
	StopBits   int    // Stop bits: 1 or 2 (default 1)
	Parity     string // Parity: N - None, E - Even, O - Odd (default E)

}
type CfgHost struct {
	Name     string
	Interval uint32
	Address  string
	CfgAddrTcp
	CfgAddrRtu
	Slave []CfgSlave
}

type CfgSlave struct {
	SlaveId    uint8
	Product    string
	Name       string
	Collection []string
}

func LoadCfgModbus(path string) (*CfgModbus, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var config CfgModbus
	if err := toml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	return &config, nil

}

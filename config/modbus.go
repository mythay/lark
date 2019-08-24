package config

import (
	"fmt"
	"reflect"
)

type Verifier interface {
	Verify() (bool, error)
}

type CfgModbus struct {
	Type    string
	Profile []CfgProfile
	Host    []CfgHost
}
type CfgProfile struct {
	Id       string
	Name     string
	Register []CfgRegister
	Range    []CfgRange
	// Inverse  bool
}
type CfgRegister struct {
	Id        uint32
	Name      string
	Start     uint16
	Quantity  uint16
	Type      string
	Inverse   bool
	Mask      uint32
	Catalog   string
	Condition string
	Action    string
	Tag       string
	//
	cCount uint16
	cType  reflect.Type
	cKind  reflect.Kind
}

func (reg *CfgRegister) End() uint16 {
	return reg.Start + reg.Quantity - 1
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

type CfgRange struct {
	Start uint16
	End   uint16
	Fixed bool
}

func (rag *CfgRange) Count() uint16 {
	return rag.End - rag.Start + 1
}

type CfgAddrTcp struct {
	IpAddr string
	Port   uint16
}

type CfgAddrRtu struct {
	Serial   string
	Baud     int
	DataBits int    // Data bits: 5, 6, 7 or 8 (default 8)
	Parity   string // Parity: N - None, E - Even, O - Odd (default E)
	StopBits int    // Stop bits: 1 or 2 (default 1)

}
type CfgHost struct {
	Name   string
	IpAddr string
	Port   uint16
	// CfgAddrTcp
	// CfgAddrRtu
	Serial   string
	Baud     int
	DataBits int    // Data bits: 5, 6, 7 or 8 (default 8)
	Parity   string // Parity: N - None, E - Even, O - Odd (default E)
	StopBits int    // Stop bits: 1 or 2 (default 1)

	Policy CfgPolicy
	Slave  []CfgSlave
}

type CfgPolicy struct {
	Interval    uint32
	Timeout     uint32
	Retry       uint8
	Concurrency uint8
	KeepAlive   bool
}

type CfgSlave struct {
	SlaveId    uint8
	Name       string
	Collection map[string][]string
}

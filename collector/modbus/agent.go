package modbus

import (
	"sync"
	"time"

	mb "github.com/goburrow/modbus"
	cfg "github.com/mythay/lark/config"
)

type mbProfile struct {
	cfg.CfgProfile
}

func (prof *mbProfile) getReg(name []string) *cfg.CfgRegister {
	return nil
}

type mbHost struct {
	host       cfg.CfgHost
	locker     sync.Mutex
	conn       mb.Client
	tcpHandler *mb.TCPClientHandler
	rtuHandler *mb.RTUClientHandler
	slave      []mbSlave
}

type mbSlave struct {
	cfg.CfgSlave
	parent *mbHost
	req    []cfg.CfgRange
	cache  map[uint16]tsValue
	regs   map[string]*regValue
}

type mbRange struct {
	cfg.CfgRange
	resp []byte
	err  error
	ts   time.Time
}

type mbReg struct {
	*cfg.CfgRegister
}
type tsValue struct {
	value [2]byte
	ts    time.Time
}

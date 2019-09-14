package modbus

import (
	"fmt"
	"sync"
	"time"

	mb "github.com/goburrow/modbus"
	cfg "github.com/mythay/lark/config"
)

type agent struct {
	hosts  []*mbHost
	config cfg.CfgModbus
}

func newAgent(config cfg.CfgModbus) (*agent, error) {
	agnt := &agent{config: config}
	for _, cfgHost := range config.Host {
		host, err := newHost(cfgHost, config.Profile)
		if err != nil {
			return nil, err
		}
		agnt.hosts = append(agnt.hosts, host)
	}
	return agnt, nil
}

func (agent *agent) step() {
	for _, host := range agent.hosts {
		host.execute()
	}
}

type mbHost struct {
	host       cfg.CfgHost
	locker     sync.Mutex
	conn       mb.Client
	tcpHandler *mb.TCPClientHandler
	rtuHandler *mb.RTUClientHandler
	slave      []*mbSlave
}

func createConnection(client *mbHost) error {
	if len(client.host.IpAddr) != 0 { // valid ip address
		port := client.host.Port
		if port == 0 {
			port = 502
		}
		handler := mb.NewTCPClientHandler(fmt.Sprintf("%s:%d", client.host.IpAddr, port))
		handler.Timeout = 1 * time.Second
		client.conn = mb.NewClient(handler)
		client.tcpHandler = handler
	} else {
		handler := mb.NewRTUClientHandler(client.host.Serial)
		client.conn = mb.NewClient(handler)
		client.rtuHandler = handler

	}
	return nil
}

func newHost(host cfg.CfgHost, profiles []cfg.CfgProfile) (*mbHost, error) {
	client := &mbHost{host: host}

	for _, s := range host.Slave {
		slv, err := newSlave(client, s, profiles)
		if err != nil {
			return nil, err
		}
		client.slave = append(client.slave, slv)
	}

	err := createConnection(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *mbHost) ReadHoldingRegisters(slaveId uint8, address, quantity uint16) (results []byte, err error) {
	client.locker.Lock()
	defer client.locker.Unlock()
	if client.tcpHandler != nil {
		client.tcpHandler.SlaveId = slaveId
	}
	if client.rtuHandler != nil {
		client.rtuHandler.SlaveId = slaveId
	}
	client.tcpHandler.SlaveId = slaveId
	return client.conn.ReadHoldingRegisters(address, quantity)
}

func (client *mbHost) preCondition() bool {
	return true
}

func (client *mbHost) postCondition() bool {
	return true
}

func (client *mbHost) execute() bool {
	if !client.preCondition() {
		return false
	}
	for _, v := range client.slave {
		v.execute()
	}
	if !client.postCondition() {
		return false
	}
	return true
}

type mbSlave struct {
	cfg.CfgSlave
	parent *mbHost
	reader *mbCache
	regs   map[string]*mbReg
}

func createCache(slv *mbSlave, profiles []cfg.CfgProfile) error {
	var rags []*cfg.CfgRange
	var regs []*cfg.CfgRegister
	for profileName, names := range slv.Collection {
		profileFound := false
		for _, profile := range profiles {
			if profile.Name == profileName { // profile found
				profileFound = true
				for _, name := range names {
					for ireg, reg := range profile.Register {
						if reg.Name == name {
							regs = append(regs, &profile.Register[ireg])
						}
					}
				}
				for irag, _ := range profile.Range {
					rags = append(rags, &profile.Range[irag])
				}
			}
		}
		if !profileFound {
			return fmt.Errorf("profile '%s' not found", profileName)
		}
	}
	reader, err := newCache(rags, regs, slv)
	if err != nil {
		return err
	}
	slv.reader = reader
	slv.regs = make(map[string]*mbReg)

	for _, reg := range regs {
		slv.regs[reg.Name] = &mbReg{CfgRegister: reg}
	}

	return nil
}

func newSlave(parent *mbHost, cfg cfg.CfgSlave, profiles []cfg.CfgProfile) (*mbSlave, error) {
	slv := &mbSlave{CfgSlave: cfg, parent: parent}
	err := createCache(slv, profiles)
	if err != nil {
		return nil, err
	}
	return slv, nil
}

func (slv *mbSlave) ReadHolding(address, quantity uint16) (results []byte, err error) {
	return slv.parent.ReadHoldingRegisters(slv.SlaveId, address, quantity)
}

func (slv *mbSlave) preCondition() bool {
	return true
}

func (slv *mbSlave) postCondition() bool {
	return true
}

func (slv *mbSlave) execute() bool {
	if !slv.preCondition() {
		return false
	}
	for _, v := range slv.regs {
		v.execute(slv.reader)
	}
	if !slv.postCondition() {
		return false
	}
	return true
}

package mbclient

import (
	"fmt"
	"net"
	"sync"
	"time"

	mb "github.com/goburrow/modbus"
	cfg "github.com/mythay/lark/config"
)

type mbRange struct {
	cfg.CfgRange
	resp []byte
	err  error
	ts   time.Time
}

type mbHost struct {
	host       cfg.CfgHost
	locker     sync.Mutex
	profile    map[string]cfg.CfgProfile
	conn       mb.Client
	tcpHandler *mb.TCPClientHandler
	rtuHandler *mb.RTUClientHandler
	slave      []mbSlave
}

func ParseRegister(client *mbHost) error {
	var slaves []mbSlave
	regs := make(map[string]*regValue)
	for _, cfgSlave := range client.host.Slave {
		device := client.profile[cfgSlave.Device]
		ranges, err := aggregateRange(cfgSlave.Collect, device)
		if err != nil {
			return err
		}
		for _, name := range cfgSlave.Collect {
			if cfgReg, ok := device.Register[name]; ok {
				regs[name] = &regValue{CfgRegister: cfgReg, name: name}
			} else {
				return fmt.Errorf("invalid collect %v", name)
			}

		}
		slaves = append(slaves, mbSlave{CfgSlave: cfgSlave, parent: client, req: ranges, regs: regs})
	}
	client.slave = slaves
	return nil
}

func CreateConnection(client *mbHost) error {
	if net.ParseIP(client.host.Address) != nil { // valid ip address
		port := client.host.Port
		if port == 0 {
			port = 502
		}
		handler := mb.NewTCPClientHandler(fmt.Sprintf("%s:%d", client.host.Address, port))
		handler.Timeout = 1 * time.Second
		client.conn = mb.NewClient(handler)
		client.tcpHandler = handler
	} else {
		handler := mb.NewRTUClientHandler(client.host.Address)
		client.conn = mb.NewClient(handler)
		client.rtuHandler = handler

	}
	return nil
}

func newHost(host cfg.CfgHost, device map[string]cfg.CfgDevice) (*mbHost, error) {
	client := &mbHost{host: host, profile: device}

	err := ParseRegister(client)
	if err != nil {
		return nil, err
	}
	err = CreateConnection(client)
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

func (client *mbHost) Once() {
	for _, slv := range client.slave {
		slv.RequestOnce()
		slv.ParseOnce()
	}
}

func (client *mbHost) Close() error {
	client.locker.Lock()
	defer client.locker.Unlock()
	if client.tcpHandler != nil {
		return client.tcpHandler.Close()
	}
	if client.rtuHandler != nil {
		return client.rtuHandler.Close()
	}
	return nil
}

type tsRegValue struct {
	value [2]byte
	ts    time.Time
}

type mbSlave struct {
	cfg.CfgSlave
	parent *mbHost
	req    []cfg.CfgRange
	resp   map[uint16]tsRegValue
	regs   map[string]*regValue
}

func (slv *mbSlave) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return slv.parent.ReadHoldingRegisters(slv.SlaveId, address, quantity)
}

func (slv *mbSlave) RequestOnce() (uint16, uint16) {
	var total, errCount uint16
	slv.resp = make(map[uint16]tsRegValue)

	for i, _ := range slv.req {
		base := slv.req[i].Base
		resp, err := slv.ReadHoldingRegisters(slv.req[i].Base, slv.req[i].Count)
		// slv.req[i].resp, slv.req[i].err, slv.req[i].ts = resp, err, time.Now()
		total += 1
		if err != nil {
			errCount += 1
		} else {
			for j := 0; j < len(resp); j += 2 {
				slv.resp[base+uint16(j/2)] = tsRegValue{[2]byte{resp[j], resp[j+1]}, time.Now()}
			}
		}
	}
	return total, errCount
}

func (slv *mbSlave) ParseOnce() {
	for name, _ := range slv.regs {
		slv.regs[name].ParseResponse(slv.resp)
	}
}

type mbRoot struct {
	cfg.CfgModbus
	hosts []*mbHost
}

func newRoot(cfg cfg.CfgModbus) (*mbRoot, error) {
	root := &mbRoot{CfgModbus: cfg}
	for _, cfgHost := range cfg.Host {
		host, err := newHost(cfgHost, cfg.Device)
		if err != nil {
			return nil, err
		}
		root.hosts = append(root.hosts, host)
	}
	return root, nil
}

func (root *mbRoot) Once() {
	for _, host := range root.hosts {
		host.Once()
	}
}

func (root *mbRoot) Close() (err error) {
	for _, host := range root.hosts {
		err = host.Close()
		if err != nil {
			break
		}
	}
	return
}

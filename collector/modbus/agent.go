package modbus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"

	mb "github.com/goburrow/modbus"
	cfg "github.com/mythay/lark/config"
)

type agent struct {
	hosts  []*mbHost
	config cfg.CfgModbus
	ctx    context.Context
	cancel context.CancelFunc
}

func newAgent(config cfg.CfgModbus) (*agent, error) {
	agnt := &agent{config: config}
	agnt.ctx, agnt.cancel = context.WithCancel(context.TODO())
	for _, cfgHost := range config.Host {
		host, err := newHost(agnt.ctx, cfgHost, config.Profile)
		if err != nil {
			return nil, err
		}
		agnt.hosts = append(agnt.hosts, host)
	}
	return agnt, nil
}

func (agent *agent) step() {
	for _, host := range agent.hosts {
		host.step()
	}
}

func (agent *agent) Start() error {
	glog.Info("hello")
	for _, host := range agent.hosts {
		err := host.Start()
		if err != nil {
			agent.cancel()
			break
		}
	}
	return nil
}

func (agent *agent) Stop() error {
	agent.cancel()
	return nil
}

type mbHost struct {
	host       cfg.CfgHost
	locker     sync.Mutex
	conn       mb.Client
	tcpHandler *mb.TCPClientHandler
	rtuHandler *mb.RTUClientHandler
	slave      []*mbSlave
	ctx        context.Context
	lastReadTs time.Time
	loop       uint64
}

func createConnection(client *mbHost) error {
	if len(client.host.IpAddr) != 0 { // valid ip address
		port := client.host.Port
		if port == 0 {
			port = 502
		}

		handler := mb.NewTCPClientHandler(fmt.Sprintf("%s:%d", client.host.IpAddr, port))
		handler.Timeout = time.Duration(client.host.Policy.Timeout) * time.Millisecond
		client.conn = mb.NewClient(handler)
		client.tcpHandler = handler
	} else {
		handler := mb.NewRTUClientHandler(client.host.Serial)
		handler.Timeout = time.Duration(client.host.Policy.Timeout) * time.Millisecond
		client.conn = mb.NewClient(handler)
		client.rtuHandler = handler

	}
	return nil
}

func newHost(ctx context.Context, host cfg.CfgHost, profiles []cfg.CfgProfile) (*mbHost, error) {
	client := &mbHost{host: host, ctx: ctx, lastReadTs: time.Now()}

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
	du := time.Duration(client.host.Policy.Interval)*time.Millisecond - time.Since(client.lastReadTs)
	if du > 0 { // calculate the interval duration
		time.Sleep(du)
	}
	r, err := client.conn.ReadHoldingRegisters(address, quantity)
	client.lastReadTs = time.Now()
	return r, err
}

func (client *mbHost) preCondition() bool {
	return true
}

func (client *mbHost) postCondition() bool {
	return true
}

func (client *mbHost) step() bool {
	if !client.preCondition() {
		return false
	}
	for _, v := range client.slave {
		v.step()
	}
	if !client.postCondition() {
		return false
	}

	return true
}

func (client *mbHost) Start() error {
	client.loop = 0
	go func() {
		period := time.Duration(client.host.Policy.Period) * time.Millisecond
		du := time.Duration(0)
		for {
			select {
			case <-client.ctx.Done():
				return
			case <-time.After(du):
				glog.Info("one loop start")
				tsBegin := time.Now()
				client.step()
				elapse := time.Since(tsBegin)
				if elapse < period {
					du = period - elapse
				} else {
					du = time.Duration(client.host.Policy.Interval)
				}
			}
			client.loop++
		}
	}()
	return nil
}

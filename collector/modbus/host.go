package modbus

import (
	"context"
	"fmt"
	"sync"
	"time"

	mb "github.com/goburrow/modbus"
	cfg "github.com/mythay/lark/config"
	log "github.com/sirupsen/logrus"
)

type mbHost struct {
	host cfg.CfgHost

	fnMetric MetricDispatcher
	bucket   chan struct{}

	locker     sync.Mutex
	conn       mb.Client
	tcpHandler *mb.TCPClientHandler
	rtuHandler *mb.RTUClientHandler
	slave      []*mbSlave
	ctx        context.Context
	lastReadTs time.Time
	loop       uint64
}

func createConnection(host *mbHost) error {
	if len(host.host.IpAddr) != 0 { // valid ip address
		port := host.host.Port
		if port == 0 {
			port = 502
		}

		handler := mb.NewTCPClientHandler(fmt.Sprintf("%s:%d", host.host.IpAddr, port))
		handler.Timeout = time.Duration(host.host.Policy.Timeout) * time.Millisecond
		host.conn = mb.NewClient(handler)
		host.tcpHandler = handler
	} else {
		handler := mb.NewRTUClientHandler(host.host.Serial)
		handler.Timeout = time.Duration(host.host.Policy.Timeout) * time.Millisecond
		host.conn = mb.NewClient(handler)
		host.rtuHandler = handler

	}
	return nil
}

func newHost(ctx context.Context, host cfg.CfgHost, profiles []cfg.CfgProfile) (*mbHost, error) {
	client := &mbHost{host: host, ctx: ctx, lastReadTs: time.Now(), bucket: make(chan struct{}, 1)}

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

func (host *mbHost) ReadHoldingRegisters(slaveId uint8, address, quantity uint16) (results []byte, err error) {
	host.locker.Lock()
	defer host.locker.Unlock()
	if host.tcpHandler != nil {
		host.tcpHandler.SlaveId = slaveId
	}
	if host.rtuHandler != nil {
		host.rtuHandler.SlaveId = slaveId
	}
	du := time.Duration(host.host.Policy.Interval)*time.Millisecond - time.Since(host.lastReadTs)
	if du > 0 { // calculate the interval duration
		time.Sleep(du)
	}
	r, err := host.conn.ReadHoldingRegisters(address, quantity)
	log.Infof("Read holding: %d-%d, got %v,%v", address, quantity, r, err)
	host.lastReadTs = time.Now()
	return r, err
}

func (host *mbHost) preCollect() bool {
	return true
}

func (host *mbHost) postCollect() bool {
	return true
}

func (host *mbHost) collect(fnMetric MetricDispatcher) bool {
	if !host.preCollect() {
		return false
	}
	for _, v := range host.slave {
		v.collect(fnMetric)
	}
	if !host.postCollect() {
		return false
	}

	return true
}

func (host *mbHost) dispatch() bool {
	return true
}

func (host *mbHost) Gather(fnMetric MetricDispatcher) error {

	if host.collect(fnMetric) {
		return nil
	}
	return fmt.Errorf("gather error")
}

func (host *mbHost) Start(fnMetric MetricDispatcher) error {
	host.loop = 0

	return nil
}

func (host *mbHost) Stop() {
}

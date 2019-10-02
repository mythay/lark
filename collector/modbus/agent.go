package modbus

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	cfg "github.com/mythay/lark/config"
)

type MetricDispatcher func(measurement string, fields map[string]interface{}, tags map[string]string, t ...time.Time)

type Agent struct {
	config   cfg.CfgModbus
	fnMetric MetricDispatcher
	hosts    []*mbHost
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewAgent() (*Agent, error) {
	agent := &Agent{}
	agent.ctx, agent.cancel = context.WithCancel(context.TODO())
	return agent, nil
}

func (agent *Agent) LoadConfig(path string) error {
	c, err := cfg.LoadCfgModbus(path)
	if err != nil {
		return err
	}
	agent.config = *c
	for _, cfgHost := range agent.config.Host {
		host, err := newHost(agent.ctx, cfgHost, agent.config.Profile)
		if err != nil {
			return err
		}
		agent.hosts = append(agent.hosts, host)
	}
	return nil
}

func (agent *Agent) Gather(fnMetric MetricDispatcher) error {

	var wg sync.WaitGroup
	for _, host := range agent.hosts {
		wg.Add(1)
		go func(host *mbHost) {
			defer wg.Done()
			_ = host.Gather(fnMetric)
		}(host)

	}
	wg.Wait()
	return nil
}

func (agent *Agent) Start(fnMetric MetricDispatcher) error {
	log.Info("hello")
	return nil
}

func (agent *Agent) Stop() {
	// agent.cancel()

}

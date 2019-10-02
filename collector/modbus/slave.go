package modbus

import (
	"fmt"

	cfg "github.com/mythay/lark/config"
)

type mbSlave struct {
	cfg.CfgSlave
	parent *mbHost
	cache  *mbCache
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
	slv.cache = reader
	slv.regs = make(map[string]*mbReg)

	for _, cfgReg := range regs {
		reg, err := NewReg(cfgReg)
		if err != nil {
			return err
		}
		slv.regs[cfgReg.Name] = reg
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

func (slv *mbSlave) preCollect() bool {
	return true
}

func (slv *mbSlave) postCollect() bool {
	return true
}

func (slv *mbSlave) collect(fnMetric MetricDispatcher) bool {
	if !slv.preCollect() {
		return false
	}
	for _, v := range slv.regs {
		v.collect(slv.cache)
	}
	if !slv.postCollect() {
		return false
	}
	// dispatch the data to some other place
	for _, v := range slv.regs {
		if v.Value != nil { //has valid value
			fnMetric("modbus", map[string]interface{}{v.Name: v.Value},
				v.tags,
				v.Ts)
		}
	}
	slv.cache.resetData()
	return true
}

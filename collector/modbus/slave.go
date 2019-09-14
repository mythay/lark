package modbus

import (
	"fmt"

	cfg "github.com/mythay/lark/config"
)

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

func (slv *mbSlave) step() bool {
	if !slv.preCondition() {
		return false
	}
	for _, v := range slv.regs {
		v.step(slv.reader)
	}
	if !slv.postCondition() {
		return false
	}
	slv.reader.resetData()
	return true
}

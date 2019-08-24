package mbclient

import (
	cfg "github.com/mythay/lark/config"
)

type innerRange struct {
	org   cfg.CfgRange
	calc  cfg.CfgRange
	valid bool
}

func (rag *innerRange) AdjustRange(reg *cfg.CfgRegister) bool {

	if reg.Start >= rag.org.Start && reg.End() <= rag.org.End {
		// don't adjust the range if the range is fixed
		if rag.org.Fixed {
			if rag.valid == false {
				rag.calc = rag.org
				rag.valid = true
			}
			return true
		}

		if rag.valid == false { // first time, all value are zero
			rag.calc.Start = reg.Start
			rag.calc.End = reg.End()
			rag.valid = true
		} else {
			if reg.Start < rag.calc.Start {
				rag.calc.Start = reg.Start
			}
			if reg.End() > rag.calc.End {
				rag.calc.End += reg.End()
			}
		}

		return true
	}
	return false
}

// func aggregateRange(collect []string, device cfg.CfgDevice) ([]cfg.CfgRange, error) {
// 	var regs []cfg.CfgRange
// 	rags := make([]innerRange, len(device.Range))
// 	for i, rag := range device.Range {
// 		rags[i].org = rag
// 	}
// 	for _, item := range collect {
// 		if reg, ok := device.Register[item]; ok {
// 			match := false
// 			for i, _ := range rags {
// 				if rags[i].AdjustRange(&reg) {
// 					match = true
// 					break
// 				}
// 			}
// 			if !match {
// 				regs = append(regs, cfg.CfgRange{Base: reg.Base, Count: reg.Count()})
// 			}
// 		} else {
// 			return nil, fmt.Errorf("'%s' reg not exist\n", item)
// 		}
// 	}
// 	for _, rag := range rags {
// 		if rag.calc.Base+rag.calc.Count > 0 {
// 			regs = append(regs, cfg.CfgRange{Base: rag.calc.Base, Count: rag.calc.Count})
// 		}
// 	}
// 	return regs, nil
// }

package modbus

import (
	"sort"
	"time"

	cfg "github.com/mythay/lark/config"
)

type MbReader interface {
	ReadHolding(address, quantity uint16) (results []byte, err error)
}

type tsValue struct {
	value [2]byte
	ts    time.Time
}

type tsRange struct {
	value []byte
	ts    time.Time
}

func (rag *tsRange) append(val tsValue) {
	rag.value = append(rag.value, val.value[:]...)
	if rag.ts.Before(val.ts) {
		rag.ts = val.ts
	}
}

func (rag *tsRange) quantity() uint16 {
	return uint16(len(rag.value) / 2)
}

type mbCache struct {
	data   map[uint16]tsValue
	rag    rangeSlice
	reader MbReader
}

func newCache(rags []*cfg.CfgRange, regs []*cfg.CfgRegister, reader MbReader) (*mbCache, error) {
	cache := &mbCache{reader: reader}
	for _, rag := range rags {
		cache.rag = append(cache.rag, readRange{org: *rag})
	}
	sort.Sort(cache.rag)
	for _, reg := range regs {
		if found := cache.locateRange(cfg.CfgRange{Start: reg.Start, End: reg.End()}); found != nil {
			found.adjustRange(reg)
		}

	}
	cache.resetData()
	return cache, nil
}

func (cache *mbCache) resetData() {
	cache.data = make(map[uint16]tsValue)
}

func (buf *mbCache) ReadThrough(start, quantity uint16) (tsRange, error) {
	var results tsRange
	// get the best range we can read
	readRange := cfg.CfgRange{Start: start, End: start + quantity - 1}
	if found := buf.locateRange(readRange); found != nil {
		if found.valid {
			readRange = found.calc
		} else {
			readRange = found.org
		}
	}

	// read action
	data, err := buf.reader.ReadHolding(readRange.Start, readRange.Count())
	if err != nil {
		return results, err
	}

	// refresh the cache
	for j := 0; j < len(data); j += 2 {
		buf.data[start+uint16(j/2)] = tsValue{[2]byte{data[j], data[j+1]}, time.Now()}
	}
	for addr := start; addr < start+quantity; addr++ {
		results.append(buf.data[addr])
	}
	return results, nil
}
func (buf *mbCache) Read(start, quantity uint16) (tsRange, error) {
	inCache := true
	var results tsRange
	for addr := start; addr < start+quantity; addr++ {
		if v, ok := buf.data[addr]; ok {
			results.append(v)
		} else {
			inCache = false
			break
		}
	}
	if inCache {
		return results, nil
	}
	// cache not match, need to read the real data from bus
	return buf.ReadThrough(start, quantity)

}

// binary search
func (buf *mbCache) locateRange(reg cfg.CfgRange) *readRange {
	s := buf.rag
	lo, hi := 0, len(s)-1
	for lo <= hi {
		m := (lo + hi) >> 1
		if s[m].org.End < reg.Start {
			lo = m + 1
		} else if s[m].org.Start > reg.End {
			hi = m - 1
		} else if s[m].org.Start <= reg.Start && s[m].org.End >= reg.End {
			return &s[m]
		} else if s[m].org.Start > reg.Start {
			hi = m - 1
		} else if s[m].org.End < reg.End {
			lo = m + 1
		} else if lo == hi {
			return nil
		}
	}
	return nil
}

type rangeSlice []readRange

func (s rangeSlice) Len() int           { return len(s) }
func (s rangeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s rangeSlice) Less(i, j int) bool { return s[i].org.Start < s[j].org.Start }

type readRange struct {
	org   cfg.CfgRange
	calc  cfg.CfgRange
	valid bool
}

func (rag *readRange) adjustRange(reg *cfg.CfgRegister) bool {

	if reg.Start >= rag.org.Start && reg.End() <= rag.org.End {
		// don't adjust the range if the range is fixed
		if rag.org.Fixed {
			return true
		}
		if rag.valid == false {
			rag.calc.Start = reg.Start
			rag.calc.End = reg.End()
			rag.valid = true
		} else {
			if reg.Start < rag.calc.Start {
				rag.calc.Start = reg.Start
			}
			if reg.End() > rag.calc.End {
				rag.calc.End = reg.End()
			}
		}

		return true
	}
	return false
}

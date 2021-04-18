package data

import (
	"sync"

	"github.com/tgmk/know-trade/types"
)

type period = string

type finReports struct {
	mx      sync.Mutex
	reports map[instrumentID][]*types.FinReport
}

func newFinReports() *finReports {
	return &finReports{
		mx:      sync.Mutex{},
		reports: make(map[instrumentID][]*types.FinReport),
	}
}

func (f *finReports) Get(id instrumentID) []*types.FinReport {
	f.mx.Lock()
	defer f.mx.Unlock()

	return f.reports[id]
}

func (f *finReports) Set(quote *types.FinReport) {
	f.mx.Lock()
	defer f.mx.Unlock()

	f.reports[quote.InstrumentID] = append(f.reports[quote.InstrumentID], quote)
}

package types

import "time"

type Match struct {
	*baseIncoming
	InstrumentID string
	Price        float64
	Size         float64
	Side         string
	Time         time.Time
	Taker        string
	Maker        string
	Other        map[string]interface{}
}

func (p *Match) Type() IncomingType {
	return IncomingMatch
}

func (p *Match) Match() *Match {
	return p
}

type GetPrintsRequest struct {
	InstrumentID string
	From         time.Time
	To           time.Time
}

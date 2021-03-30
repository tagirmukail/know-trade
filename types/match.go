package types

import "time"

type Match struct {
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

func (p *Match) Candle() *Candle {
	panic("not implemented")
}

func (p *Match) OrderBook() *OrderBook {
	panic("not implemented")
}

func (p *Match) Match() *Match {
	return p
}

func (p *Match) Order() *Order {
	panic("not implemented")
}

type GetPrintsRequest struct {
	InstrumentID string
	From         time.Time
	To           time.Time
}

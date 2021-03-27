package types

import "time"

type Print struct {
	Symbol string
	Price  float64
	Size   float64
	Side   string
	Time   time.Time
	Other  map[string]interface{}
}

func (p *Print) Type() IncomingType {
	return IncomingPrint
}

func (p *Print) Candle() *Candle {
	panic("not implemented")
}

func (p *Print) OrderBook() *OrderBook {
	panic("not implemented")
}

func (p *Print) Print() *Print {
	return p
}

func (p *Print) Order() *Order {
	panic("not implemented")
}

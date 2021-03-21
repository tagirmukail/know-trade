package types

import "time"

type Candle struct {
	Close  float64
	Open   float64
	High   float64
	Low    float64
	Time   time.Time
	Volume float64
	Other  map[string]interface{}
}

func (c *Candle) Type() IncomingType {
	return IncomingCandle
}

func (c *Candle) Candle() *Candle {
	return c
}

func (c *Candle) OrderBook() *OrderBook {
	panic("not implemented")
}

func (c *Candle) Print() *Print {
	panic("not implemented")
}

package types

import "time"

type Period string

const (
	Period1Min  = "1min"
	Period3Min  = "3min"
	Period5Min  = "5min"
	Period10Min = "10min"
	Period15Min = "15min"
	Period30Min = "30min"
	Period1Hour = "1hour"
	Period2Hour = "2hour"
	Period4Hour = "4hour"
	Period1Day  = "1day"
)

type Candle struct {
	Close  float64
	Open   float64
	High   float64
	Low    float64
	Time   time.Time
	Volume float64
	Period Period
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

func (c *Candle) Order() *Order {
	panic("not implemented")
}

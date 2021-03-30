package types

import "time"

type OrderBook struct {
	Price        float64
	Size         float64
	Time         time.Time
	InstrumentID string
	Other        map[string]interface{}
}

func (o *OrderBook) Type() IncomingType {
	return IncomingOrderBook
}

func (o *OrderBook) Candle() *Candle {
	panic("not implemented")
}

func (o *OrderBook) OrderBook() *OrderBook {
	return o
}

func (o *OrderBook) Match() *Match {
	panic("not implemented")
}

func (o *OrderBook) Order() *Order {
	panic("not implemented")
}

type GetOrderBookRequest struct {
	InstrumentID string
	Dept         int
}

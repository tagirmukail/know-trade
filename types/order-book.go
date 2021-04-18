package types

import "time"

type OrderBook struct {
	*baseIncoming
	Price        float64
	Size         float64
	Time         time.Time
	InstrumentID string
	Other        map[string]interface{}
}

func (o *OrderBook) Type() IncomingType {
	return IncomingOrderBook
}

func (o *OrderBook) OrderBook() *OrderBook {
	return o
}

type GetOrderBookRequest struct {
	InstrumentID string
	Dept         int
}

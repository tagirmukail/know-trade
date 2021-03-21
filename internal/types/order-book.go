package types

import "time"

type OrderBook struct {
	Price float64
	Size  float64
	Time  time.Time
	Other map[string]interface{}
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

func (o *OrderBook) Print() *Print {
	panic("not implemented")
}

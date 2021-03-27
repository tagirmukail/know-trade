package types

import "time"

type Status uint8

const (
	_ Status = iota
	Opened
	Filled
	PartialFilled
	Canceled
	Failed
)

func (o *Order) Type() IncomingType {
	return IncomingOrder
}

type Order struct {
	ID     string
	Symbol string
	Side   string
	Status Status
	Price  float64
	Size   float64
	Time   time.Time
	Other  map[string]interface{}
}

func (o *Order) Order() *Order {
	return o
}

func (o *Order) Candle() *Candle {
	panic("not implemented")
}

func (o *Order) OrderBook() *OrderBook {
	panic("not implemented")
}

func (o *Order) Print() *Print {
	panic("not implemented")
}

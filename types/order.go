package types

import "time"

type Status string

const (
	Open          Status = "open"
	Filled        Status = "filled"
	PartialFilled Status = "partial_filled"
	Canceled      Status = "canceled"
	Failed        Status = "failed"
)

type Order struct {
	*baseIncoming
	ID           string
	InstrumentID string
	Side         string
	Status       Status
	Price        float64
	Size         float64
	Time         time.Time
	Other        map[string]interface{}
}

func (o *Order) Type() IncomingType {
	return IncomingOrder
}

func (o *Order) Order() *Order {
	return o
}

type LimitOrderRequest struct {
	Price        float64
	Size         float64
	InstrumentID string
	Side         string
}

type MarketOrderRequest struct {
	Price        float64
	Size         float64
	InstrumentID string
	Side         string
}

type CancelOrderRequest struct {
	OrderID string
}

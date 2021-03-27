package types

type IncomingType uint8

const (
	_ IncomingType = iota
	IncomingCandle
	IncomingOrderBook
	IncomingPrint
	IncomingOrder
)

type Incoming interface {
	Type() IncomingType
	Candle() *Candle
	OrderBook() *OrderBook
	Print() *Print
	Order() *Order
}

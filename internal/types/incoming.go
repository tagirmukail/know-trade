package types

type IncomingType uint8

const (
	_ IncomingType = iota
	IncomingCandle
	IncomingOrderBook
	IncomingMatch
	IncomingOrder
)

type Incoming interface {
	Type() IncomingType
	Candle() *Candle
	OrderBook() *OrderBook
	Match() *Match
	Order() *Order
}

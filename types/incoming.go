package types

type IncomingType uint8

const (
	_ IncomingType = iota
	IncomingCandle
	IncomingOrderBook
	IncomingMatch
	IncomingOrder
	IncomingFinReport
)

type Incoming interface {
	Type() IncomingType
	Candle() *Candle
	OrderBook() *OrderBook
	Match() *Match
	Order() *Order
	FinReport() *FinReport
}

type baseIncoming struct{}

func (b *baseIncoming) Type() IncomingType {
	panic("not implemented")
}
func (b *baseIncoming) Candle() *Candle {
	panic("not implemented")
}

func (b *baseIncoming) OrderBook() *OrderBook {
	panic("not implemented")
}

func (b *baseIncoming) Match() *Match {
	panic("not implemented")
}

func (b *baseIncoming) Order() *Order {
	panic("not implemented")
}

func (b *baseIncoming) FinReport() *FinReport {
	panic("not implemented")
}

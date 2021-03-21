package testcli

import (
	"log"
	"strconv"
	"time"
)

type MatchMessage struct {
	Sequence     string `json:"sequence"`
	Type         string `json:"type"`
	Symbol       string `json:"symbol"`
	Side         string `json:"side"`
	Price        string `json:"price"`
	Size         string `json:"size"`
	TradeID      string `json:"tradeId"`
	TakerOrderID string `json:"takerOrderId"`
	MakerOrderID string `json:"makerOrderId"`
	Time         string `json:"time"`
}

func (m *MatchMessage) GetType() string {
	return m.Type
}

func (m *MatchMessage) GetSide() string {
	return m.Side
}

func (m *MatchMessage) GetSymbol() string {
	return m.Symbol
}

func (m *MatchMessage) GetPrice() float64 {
	p, err := strconv.ParseFloat(m.Price, 64)
	if err != nil {
		log.Printf("kukocin MatchMessage  get price error: %v", err)
		return 0
	}

	return p
}

func (m *MatchMessage) GetSize() float64 {
	sz, err := strconv.ParseFloat(m.Size, 64)
	if err != nil {
		log.Printf("kukocin MatchMessage  get size error: %v", err)
		return 0
	}

	return sz
}

func (m *MatchMessage) GetTime() time.Time {
	nsec, err := strconv.ParseInt(m.Time, 10, 64)
	if err != nil {
		log.Printf("kukocin MatchMessage  get time error: %v", err)
		return time.Time{}
	}

	return time.Unix(0, nsec)
}

package config

import (
	"time"

	"github.com/tgmk/know-trade/internal/executor"
)

type RunType uint8

const (
	TickerRun RunType = iota
	EveryCandleRun
	EveryPrintRun
	ByOthersRun
)

type Data struct {
	CandlesSize   int
	OrderBookSize int
	PrintsSize    int
	Other         map[string]interface{}
}

type Run struct {
	HowRun         RunType
	TickerInterval time.Duration
	Other          map[string]interface{}
}

type Test struct {
	Other map[string]interface{}
}

type Config struct {
	Run
	Data
	Test
	OrderExecutor executor.Order
	Other         map[string]interface{}
}

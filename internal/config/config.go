package config

import (
	"time"
)

type RunType uint8

const (
	TickerRun RunType = iota
	EveryCandleRun
	EveryPrintRun
)

type Data struct {
	CandlesSize   int
	OrderBookSize int
	PrintsSize    int
	Other         map[string]interface{}
}

type Run struct {
	Symbol         string
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
	Other map[string]interface{}
}

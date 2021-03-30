package config

import (
	"time"
)

type RunType uint8

const (
	TickerRun RunType = iota
	EveryCandleRun
	EveryMatchRun
)

type Data struct {
	CandlesSize   int
	OrderBookSize int
	MatchesSize   int
	Other         map[string]interface{}
}

type Run struct {
	InstrumentID   string
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

package config

type RunType uint8

const (
	TickerRun RunType = iota
	EveryCandleRun
	EveryMatchRun
	EveryPositionChangeRun
	EveryFinReport
)

type Data struct {
	CandlesSize   int
	OrderBookSize int
	MatchesSize   int
	Other         map[string]interface{}
}

type Test struct {
	Other map[string]interface{}
}

type Config struct {
	Data
	Test
	Other map[string]interface{}
}

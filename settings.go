package knowtrade

import (
	"time"

	"github.com/tgmk/know-trade/config"
)

// RunSettings trade strategy settings
type RunSettings struct {
	config.RunType
	Handler
	InstrumentID   string
	TickerInterval time.Duration
	Other          map[string]interface{}
}

type HowRun map[config.RunType]*RunSettings

func NewHowRun(runs ...RunSettings) HowRun {
	hr := map[config.RunType]*RunSettings{}

	for _, pair := range runs {
		hr[pair.RunType] = &pair
	}

	return hr
}

func (hr HowRun) GetRunTypes() (result []config.RunType) {
	for rt, _ := range hr {
		result = append(result, rt)
	}

	return result
}

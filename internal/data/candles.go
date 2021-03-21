package data

import (
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type candles struct {
	sync.Mutex
	candles []*types.Candle

	size int
}

func newCandles(size int) *candles {
	return &candles{
		Mutex:   sync.Mutex{},
		candles: make([]*types.Candle, 0),
		size:    size,
	}
}

func (c *candles) Set(candle *types.Candle) {
	c.Lock()
	defer c.Unlock()

	first := len(c.candles) - c.size
	if first < 0 {
		first = 0
	}

	c.candles = append(c.candles[first:], candle)
}

func (c *candles) GetLast() *types.Candle {
	c.Lock()
	defer c.Unlock()

	return c.candles[len(c.candles)-1]
}

func (c *candles) Get() []*types.Candle {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.Candle, len(c.candles))
	copy(res, c.candles)

	return res
}

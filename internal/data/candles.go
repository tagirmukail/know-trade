package data

import (
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type candles struct {
	sync.Mutex
	candles map[types.Period][]*types.Candle

	size int
}

func newCandles(size int) *candles {
	return &candles{
		Mutex:   sync.Mutex{},
		candles: make(map[types.Period][]*types.Candle),
		size:    size,
	}
}

func (c *candles) Set(candle *types.Candle) {
	c.Lock()
	defer c.Unlock()

	_, ok := c.candles[candle.Period]
	if !ok {
		c.candles[candle.Period] = make([]*types.Candle, 0)
	}

	first := len(c.candles[candle.Period]) - c.size
	if first < 0 {
		first = 0
	}

	c.candles[candle.Period] = append(c.candles[candle.Period][first:], candle)
}

func (c *candles) GetLast(period types.Period) *types.Candle {
	c.Lock()
	defer c.Unlock()

	return c.candles[period][len(c.candles[period])-1]
}

func (c *candles) Get(period types.Period) []*types.Candle {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.Candle, len(c.candles[period]))
	copy(res, c.candles[period])

	return res
}

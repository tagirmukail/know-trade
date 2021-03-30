package data

import (
	"sync"

	"github.com/tgmk/know-trade/types"
)

type candles struct {
	sync.Mutex
	candles map[instrumentID]map[types.Period][]*types.Candle

	size int
}

func newCandles(size int) *candles {
	return &candles{
		Mutex:   sync.Mutex{},
		candles: make(map[instrumentID]map[types.Period][]*types.Candle),
		size:    size,
	}
}

func (c *candles) Set(candle *types.Candle) {
	c.Lock()
	defer c.Unlock()

	_, ok := c.candles[candle.InstrumentID]
	if !ok {
		c.candles[candle.InstrumentID] = map[types.Period][]*types.Candle{}
	}

	_, ok = c.candles[candle.InstrumentID][candle.Period]
	if !ok {
		c.candles[candle.InstrumentID][candle.Period] = make([]*types.Candle, 0)
	}

	first := len(c.candles[candle.InstrumentID][candle.Period]) - c.size
	if first < 0 {
		first = 0
	}

	c.candles[candle.InstrumentID][candle.Period] = append(c.candles[candle.InstrumentID][candle.Period][first:], candle)
}

func (c *candles) GetLast(instr instrumentID, period types.Period) *types.Candle {
	c.Lock()
	defer c.Unlock()

	if len(c.candles[instr]) == 0 {
		return &types.Candle{}
	}

	if len(c.candles[instr][period]) == 0 {
		return &types.Candle{}
	}

	return c.candles[instr][period][len(c.candles[instr][period])-1]
}

func (c *candles) Get(instr instrumentID, period types.Period) (res []*types.Candle) {
	c.Lock()
	defer c.Unlock()

	if len(c.candles[instr]) == 0 {
		return res
	}

	res = make([]*types.Candle, len(c.candles[instr][period]))
	copy(res, c.candles[instr][period])

	return res
}

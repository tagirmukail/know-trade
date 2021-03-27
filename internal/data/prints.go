package data

import (
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type prints struct {
	sync.Mutex
	cache map[symbol][]*types.Print

	size int
}

func newPrints(size int) *prints {
	return &prints{
		Mutex: sync.Mutex{},
		cache: make(map[symbol][]*types.Print),
		size:  size,
	}
}

func (c *prints) Set(print *types.Print) {
	c.Lock()
	defer c.Unlock()

	first := len(c.cache[print.Symbol]) - c.size
	if first < 0 {
		first = 0
	}

	c.cache[print.Symbol] = append(c.cache[print.Symbol][first:], print)
}

func (c *prints) GetLast(s symbol) *types.Print {
	c.Lock()
	defer c.Unlock()

	return c.cache[s][len(c.cache)-1]
}

func (c *prints) Get(s symbol) []*types.Print {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.Print, len(c.cache[s]))
	copy(res, c.cache[s])

	return res
}

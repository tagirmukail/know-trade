package data

import (
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type prints struct {
	sync.Mutex
	cache []*types.Print

	size int
}

func newPrints(size int) *prints {
	return &prints{
		Mutex: sync.Mutex{},
		cache: make([]*types.Print, 0),
		size:  size,
	}
}

func (c *prints) Set(print *types.Print) {
	c.Lock()
	defer c.Unlock()

	first := len(c.cache) - c.size
	if first < 0 {
		first = 0
	}

	c.cache = append(c.cache[first:], print)
}

func (c *prints) GetLast() *types.Print {
	c.Lock()
	defer c.Unlock()

	return c.cache[len(c.cache)-1]
}

func (c *prints) Get() []*types.Print {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.Print, len(c.cache))
	copy(res, c.cache)

	return res
}

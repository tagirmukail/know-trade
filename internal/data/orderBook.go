package data

import (
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type orderBookCache struct {
	sync.Mutex
	cache []*types.OrderBook

	size int
}

func newOrderBookCache(size int) *orderBookCache {
	return &orderBookCache{
		Mutex: sync.Mutex{},
		cache: make([]*types.OrderBook, 0),
		size:  size,
	}
}

func (c *orderBookCache) Set(orderBook *types.OrderBook) {
	c.Lock()
	defer c.Unlock()

	first := len(c.cache) - c.size
	if first < 0 {
		first = 0
	}

	c.cache = append(c.cache[first:], orderBook)
}

func (c *orderBookCache) GetLast() *types.OrderBook {
	c.Lock()
	defer c.Unlock()

	return c.cache[len(c.cache)-1]
}

func (c *orderBookCache) Get() []*types.OrderBook {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.OrderBook, len(c.cache))
	copy(res, c.cache)

	return res
}

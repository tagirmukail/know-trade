package data

import (
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type symbol = string

type orderBookCache struct {
	sync.Mutex
	cache map[symbol][]*types.OrderBook

	size int
}

func newOrderBookCache(size int) *orderBookCache {
	return &orderBookCache{
		Mutex: sync.Mutex{},
		cache: make(map[symbol][]*types.OrderBook),
		size:  size,
	}
}

func (c *orderBookCache) Set(orderBook *types.OrderBook) {
	c.Lock()
	defer c.Unlock()

	first := len(c.cache[orderBook.Symbol]) - c.size
	if first < 0 {
		first = 0
	}

	c.cache[orderBook.Symbol] = append(c.cache[orderBook.Symbol][first:], orderBook)
}

func (c *orderBookCache) GetLast(sym string) *types.OrderBook {
	c.Lock()
	defer c.Unlock()

	return c.cache[sym][len(c.cache)-1]
}

func (c *orderBookCache) Get(sym string) []*types.OrderBook {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.OrderBook, len(c.cache[sym]))
	copy(res, c.cache[sym])

	return res
}

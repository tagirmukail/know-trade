package data

import (
	"sync"

	"github.com/tgmk/know-trade/types"
)

type instrumentID = string

type orderBookCache struct {
	sync.Mutex
	cache map[instrumentID][]*types.OrderBook

	size int
}

func newOrderBookCache(size int) *orderBookCache {
	return &orderBookCache{
		Mutex: sync.Mutex{},
		cache: make(map[instrumentID][]*types.OrderBook),
		size:  size,
	}
}

func (c *orderBookCache) Set(orderBook *types.OrderBook) {
	c.Lock()
	defer c.Unlock()

	first := len(c.cache[orderBook.InstrumentID]) - c.size
	if first < 0 {
		first = 0
	}

	c.cache[orderBook.InstrumentID] = append(c.cache[orderBook.InstrumentID][first:], orderBook)
}

func (c *orderBookCache) GetLast(instrumentID string) *types.OrderBook {
	c.Lock()
	defer c.Unlock()

	if len(c.cache[instrumentID]) == 0 {
		return &types.OrderBook{}
	}

	return c.cache[instrumentID][len(c.cache)-1]
}

func (c *orderBookCache) Get(sym string) []*types.OrderBook {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.OrderBook, len(c.cache[sym]))
	copy(res, c.cache[sym])

	return res
}

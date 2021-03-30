package data

import (
	"sync"

	"github.com/tgmk/know-trade/types"
)

type matches struct {
	sync.Mutex
	cache map[instrumentID][]*types.Match

	size int
}

func newMatches(size int) *matches {
	return &matches{
		Mutex: sync.Mutex{},
		cache: make(map[instrumentID][]*types.Match),
		size:  size,
	}
}

func (c *matches) Set(match *types.Match) {
	c.Lock()
	defer c.Unlock()

	first := len(c.cache[match.InstrumentID]) - c.size
	if first < 0 {
		first = 0
	}

	c.cache[match.InstrumentID] = append(c.cache[match.InstrumentID][first:], match)
}

func (c *matches) GetLast(s instrumentID) *types.Match {
	c.Lock()
	defer c.Unlock()

	if len(c.cache[s]) == 0 {
		return &types.Match{}
	}

	return c.cache[s][len(c.cache)-1]
}

func (c *matches) Get(s instrumentID) []*types.Match {
	c.Lock()
	defer c.Unlock()

	res := make([]*types.Match, len(c.cache[s]))
	copy(res, c.cache[s])

	return res
}

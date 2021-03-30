package data

import (
	"sort"
	"sync"

	"github.com/tgmk/know-trade/types"
)

type id = string

type position struct {
	sync.Mutex
	cache map[id]*types.Order
}

func newPosition() *position {
	return &position{
		Mutex: sync.Mutex{},
		cache: make(map[id]*types.Order),
	}
}

func (o *position) Set(order *types.Order) {
	o.Lock()
	defer o.Unlock()

	o.cache[order.ID] = order
}

func (o *position) List() []*types.Order {
	o.Lock()
	defer o.Unlock()

	result := make([]*types.Order, len(o.cache))
	for _, order := range o.cache {
		result = append(result, order)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result
}

func (o *position) Get(id id) *types.Order {
	o.Lock()
	defer o.Unlock()

	order, ok := o.cache[id]
	if !ok {
		return nil
	}

	return order
}

// Update remove canceled and filled position from cache
func (o *position) Update() {
	o.Lock()
	defer o.Unlock()

	for id, order := range o.cache {
		if order.Status != types.Filled && order.Status != types.Canceled {
			continue
		}

		delete(o.cache, id)
	}
}

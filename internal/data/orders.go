package data

import (
	"sort"
	"sync"

	"github.com/tgmk/know-trade/internal/types"
)

type id = string

type orders struct {
	sync.Mutex
	cache map[id]*types.Order
}

func newOrders() *orders {
	return &orders{
		Mutex: sync.Mutex{},
		cache: make(map[id]*types.Order),
	}
}

func (o *orders) Set(order *types.Order) {
	o.Lock()
	defer o.Unlock()

	o.cache[order.ID] = order
}

func (o *orders) List() []*types.Order {
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

func (o *orders) Get(id id) *types.Order {
	order, ok := o.cache[id]
	if !ok {
		return nil
	}

	return order
}

// Update remove canceled and filled orders from cache
func (o *orders) Update() {
	o.Lock()
	defer o.Unlock()

	for id, order := range o.cache {
		if order.Status != types.Filled && order.Status != types.Canceled {
			continue
		}

		delete(o.cache, id)
	}
}

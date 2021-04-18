package sync

import "sync"

type ChMap struct {
	mx sync.Mutex
	m  map[string]chan struct{}
}

func New() *ChMap {
	return &ChMap{
		mx: sync.Mutex{},
		m:  make(map[string]chan struct{}),
	}
}

func (c *ChMap) Set(instrument string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	_, ok := c.m[instrument]
	if ok {
		return
	}

	c.m[instrument] = make(chan struct{}, 1024)
}

func (c *ChMap) Upsert(instrument string) chan struct{} {
	c.mx.Lock()
	defer c.mx.Unlock()

	ch, ok := c.m[instrument]
	if ok {
		return ch
	}

	c.m[instrument] = make(chan struct{}, 1024)

	return c.m[instrument]
}

func (c *ChMap) Get(instrument string) (chan struct{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	ch, ok := c.m[instrument]

	return ch, ok
}

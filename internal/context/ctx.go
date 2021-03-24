package context

import (
	"context"

	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/data"
)

type Context struct {
	context.Context

	cfg *config.Config

	d *data.Data
}

func New(ctx context.Context, cfg *config.Config) *Context {
	return &Context{
		Context: ctx,
		cfg:     cfg,
		d:       data.New(ctx, cfg),
	}
}

func (c *Context) SetConfig(cfg *config.Config) {
	if c == nil {
		panic("Context is not initialized")
	}

	c.cfg = cfg
}

func (c *Context) GetConfig() *config.Config {
	if c == nil {
		panic("Context is not initialized")
	}

	return c.cfg
}

func (c *Context) SetData(d *data.Data) {
	if c == nil {
		panic("Context is not initialized")
	}

	c.d = d
}

func (c *Context) GetData() *data.Data {
	if c == nil {
		panic("Context is not initialized")
	}

	return c.d
}

func (c *Context) SetCtx(ctx context.Context) {
	if c == nil {
		panic("Context is not initialized")
	}

	c.Context = ctx
}

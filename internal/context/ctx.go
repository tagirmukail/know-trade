package context

import (
	"context"

	"github.com/tgmk/know-trade/internal/exchange"

	"github.com/tgmk/know-trade/internal/db"

	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/data"
)

type Context struct {
	context.Context

	cfg *config.Config

	d *data.Data

	db db.IDB

	exchangeClient exchange.IClient
}

func New(ctx context.Context, cfg *config.Config) *Context {
	return &Context{
		Context: ctx,
		cfg:     cfg,
		d:       data.New(ctx, cfg),
	}
}

func (c *Context) SetDB(db db.IDB) {
	if c == nil {
		panic("Context is not initialized")
	}

	c.db = db
}

func (c *Context) GetDB() db.IDB {
	if c == nil {
		panic("Context is not initialized")
	}

	return c.db
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

func (c *Context) SetExchangeClient(client exchange.IClient) {
	c.exchangeClient = client
}

func (c *Context) GetExchangeClient() exchange.IClient {
	return c.exchangeClient
}

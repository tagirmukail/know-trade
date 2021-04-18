package context

import (
	"context"

	"github.com/tgmk/know-trade/db"

	"github.com/tgmk/know-trade/config"
	"github.com/tgmk/know-trade/data"
)

type Context struct {
	context.Context

	context.CancelFunc

	cfg *config.Config

	d *data.Data

	db db.IDB
}

func New(cfg *config.Config, runTypes map[string]map[config.RunType]struct{}) *Context {
	if len(runTypes) == 0 {
		panic("run settings not configured")
	}

	ctx := &Context{
		cfg: cfg,
	}

	ctx.Context, ctx.CancelFunc = context.WithCancel(context.Background())

	ctx.d = data.New(ctx.Context, cfg, runTypes)

	return ctx
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

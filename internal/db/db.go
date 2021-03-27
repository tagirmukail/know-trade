package db

import (
	ctx "github.com/tgmk/know-trade/internal/context"
)

type CollectionName string

type IDB interface {
	Collection(name CollectionName) Collection
}

type Collection interface {
	Save(ctx *ctx.Context, in interface{}) error
	Get(ctx *ctx.Context) interface{}
	List(ctx *ctx.Context) []interface{}
	Remove(ctx *ctx.Context, in interface{}) error
}

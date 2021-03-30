package db

import (
	"context"
)

type CollectionName string

type IDB interface {
	Collection(name CollectionName) Collection
}

type Collection interface {
	Save(ctx context.Context, in interface{}) error
	New(ctx context.Context, in interface{}) error
	Update(ctx context.Context, in interface{}) error
	Get(ctx context.Context) interface{}
	List(ctx context.Context) []interface{}
	Remove(ctx context.Context, in interface{}) error
}

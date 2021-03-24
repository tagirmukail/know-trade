package executor

import (
	"context"

	"github.com/tgmk/know-trade/internal/types"
)

type ExchangeClient interface {
	Market(ctx context.Context, symbol, side string, size float64) (*types.Order, error)
	Limit(ctx context.Context, symbol, side string, price, size float64) (*types.Order, error)
	Cancel(ctx context.Context, orderID string) (*types.Order, error)
}

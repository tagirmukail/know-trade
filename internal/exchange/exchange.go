package exchange

import (
	"context"

	"github.com/tgmk/know-trade/internal/types"
)

type IClient interface {
	// Order
	Market(ctx context.Context, symbol, side string, size float64) (*types.Order, error)
	Limit(ctx context.Context, symbol, side string, price, size float64) (*types.Order, error)
	Cancel(ctx context.Context, orderID string) (*types.Order, error)
	// Candles
	GetCandles(ctx context.Context) ([]*types.Candle, error)
	GetOrderBook(ctx context.Context) ([]*types.OrderBook, error)
	GetPrints(ctx context.Context) ([]*types.Print, error)
}

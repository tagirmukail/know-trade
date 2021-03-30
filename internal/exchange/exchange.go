package exchange

import (
	"context"

	"github.com/tgmk/know-trade/internal/types"
)

type IClient interface {
	// Order
	Market(ctx context.Context, req *types.MarketOrderRequest) (*types.Order, error)
	Limit(ctx context.Context, req *types.LimitOrderRequest) (*types.Order, error)
	Cancel(ctx context.Context, req *types.CancelOrderRequest) (*types.Order, error)
	// Candles
	GetCandles(ctx context.Context, req *types.GetCandlesRequest) ([]*types.Candle, error)
	GetOrderBook(ctx context.Context, req *types.GetOrderBookRequest) ([]*types.OrderBook, error)
	GetPrints(ctx context.Context, req *types.GetPrintsRequest) ([]*types.Match, error)
}

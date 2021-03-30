package testcli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/tgmk/know-trade/types"
)

type TestExchangeClient struct {
	ctx context.Context

	fee        float64
	feePercent float64
	earning    float64
	balance    float64

	mx             sync.Mutex
	executedOrders map[string]*types.Order
	canceledOrders map[string]*types.Order
}

func New(ctx context.Context, feeFixed, feePercent, balance float64) *TestExchangeClient {
	return &TestExchangeClient{
		ctx: ctx,

		fee:            feeFixed,
		feePercent:     feePercent,
		balance:        balance,
		mx:             sync.Mutex{},
		executedOrders: make(map[string]*types.Order),
		canceledOrders: make(map[string]*types.Order),
	}
}

func (c *TestExchangeClient) Market(ctx context.Context, req *types.MarketOrderRequest) (*types.Order, error) {
	panic("no use")
}

func (c *TestExchangeClient) Limit(_ context.Context, req *types.LimitOrderRequest) (*types.Order, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if strings.TrimSpace(req.InstrumentID) == "" {
		return nil, errors.New("invalid instrumentID")
	}

	side := strings.ToLower(req.Side)
	if side != "sell" && side != "buy" {
		return nil, errors.New("invalid side")
	}

	if req.Price <= 0 {
		return nil, errors.New("invalid price")
	}

	if req.Size <= 0 {
		return nil, errors.New("invalid size")
	}

	err := c.calculateEarning(side, req.Price, req.Size)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	o := &types.Order{
		ID:           id.String(),
		InstrumentID: req.InstrumentID,
		Side:         side,
		Status:       types.Filled,
		Price:        req.Price,
		Size:         req.Size,
		Other: map[string]interface{}{
			"tag": "backtest",
		},
	}

	c.executedOrders[id.String()] = o

	return o, nil
}

func (c *TestExchangeClient) Cancel(_ context.Context, req *types.CancelOrderRequest) (*types.Order, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	o, ok := c.executedOrders[req.OrderID]
	if !ok {
		return nil, fmt.Errorf("order: %v does not exist", req.OrderID)
	}

	o.Status = types.Canceled

	c.canceledOrders[req.OrderID] = o

	delete(c.executedOrders, req.OrderID)

	c.calculateCanceled(o.Side, o.Price, o.Size)

	return o, nil
}

func (c *TestExchangeClient) GetCandles(ctx context.Context, req *types.GetCandlesRequest) ([]*types.Candle, error) {
	panic("not implemented")
}

func (c *TestExchangeClient) GetOrderBook(ctx context.Context, req *types.GetOrderBookRequest) ([]*types.OrderBook, error) {
	panic("not implemented")
}

func (c *TestExchangeClient) GetPrints(ctx context.Context, req *types.GetPrintsRequest) ([]*types.Match, error) {
	panic("not implemented")
}

func (c *TestExchangeClient) calculateCanceled(side string, price, size float64) {
	orderPrice := size * price
	var orderFee float64
	if c.fee > 0 {
		orderFee = size * c.fee
	} else if c.feePercent > 0 {
		orderFee = orderPrice * c.feePercent / 100
	}

	switch strings.ToLower(side) {
	case "sell":
		c.balance += orderFee
		c.earning += orderFee
		c.earning -= orderPrice
		c.balance -= orderPrice
	case "buy":
		c.balance += orderFee
		c.earning += orderFee
		c.earning += orderPrice
		c.balance += orderPrice
	default:
	}

	return
}

func (c *TestExchangeClient) calculateEarning(side string, price, size float64) error {
	orderPrice := size * price

	var orderFee float64
	if c.fee > 0 {
		orderFee = size * c.fee
	} else if c.feePercent > 0 {
		orderFee = orderPrice * c.feePercent / 100
	}

	if c.balance < orderPrice {
		return fmt.Errorf("not enough funds on balance: %v, deal price: %v", c.balance, orderPrice)
	}

	switch strings.ToLower(side) {
	case "sell":
		c.balance -= orderFee
		c.earning -= orderFee
		c.earning += orderPrice
		c.balance += orderPrice
	case "buy":
		c.balance -= orderFee
		c.earning -= orderFee
		c.earning -= orderPrice
		c.balance -= orderPrice
	default:
	}

	return nil
}

func (c *TestExchangeClient) Result() (resp struct {
	Earning, Balance float64
}) {
	resp.Earning = c.earning
	resp.Balance = c.balance

	return
}

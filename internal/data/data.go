package data

import (
	"context"

	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/types"
)

type Data struct {
	ctx context.Context

	config *config.Config

	incomingDataCh   chan types.Incoming
	candleReceivedCh chan struct{}
	printCh          chan struct{}

	candles        *candles
	orderBookCache *orderBookCache
	printsCache    *prints
	orders         *orders
}

func New(ctx context.Context, cfg *config.Config) *Data {
	return &Data{
		ctx: ctx,

		config: cfg,

		incomingDataCh:   make(chan types.Incoming, 2024),
		candleReceivedCh: make(chan struct{}, 1024),
		printCh:          make(chan struct{}, 2024),

		candles:        newCandles(cfg.Data.CandlesSize),
		orderBookCache: newOrderBookCache(cfg.Data.OrderBookSize),
		printsCache:    newPrints(cfg.Data.PrintsSize),
		orders:         newOrders(),
	}
}

func (d *Data) GetCandles() *candles {
	return d.candles
}

func (d *Data) GetOrderBookCache() *orderBookCache {
	return d.orderBookCache
}

func (d *Data) GetPrints() *prints {
	return d.printsCache
}

func (d *Data) SendToIncomingCh(inc types.Incoming) {
	d.incomingDataCh <- inc
}

func (d *Data) SendToCandlesCh() {
	d.candleReceivedCh <- struct{}{}
}

func (d *Data) SendToPrintCh() {
	d.printCh <- struct{}{}
}

func (d *Data) IncomingCh() chan types.Incoming {
	return d.incomingDataCh
}

func (d *Data) CandleCh() chan struct{} {
	return d.candleReceivedCh
}

func (d *Data) PrintCh() chan struct{} {
	return d.printCh
}

func (d *Data) Process() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case inc := <-d.incomingDataCh:
			switch inc.Type() {
			case types.IncomingCandle:
				candle := inc.Candle()

				if d.config.HowRun == config.EveryCandleRun {
					d.SendToCandlesCh()
				}

				d.candles.Set(candle)
			case types.IncomingOrderBook:
				ob := inc.OrderBook()

				d.orderBookCache.Set(ob)
			case types.IncomingPrint:
				p := inc.Print()

				if d.config.HowRun == config.EveryPrintRun {
					d.SendToPrintCh()
				}

				d.printsCache.Set(p)
			case types.IncomingOrder:
				o := inc.Order()

				d.orders.Set(o)
			default:
				panic("unknown incoming data")
			}
		}
	}
}

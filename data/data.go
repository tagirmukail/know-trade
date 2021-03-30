package data

import (
	"context"

	"github.com/tgmk/know-trade/config"
	"github.com/tgmk/know-trade/types"
)

type Data struct {
	ctx context.Context

	config *config.Config

	runTypes map[config.RunType]struct{}

	incomingDataCh   chan types.Incoming
	candleReceivedCh chan struct{}
	matchCh          chan struct{}
	positionCh       chan struct{}

	candles        *candles
	orderBookCache *orderBookCache
	matchesCache   *matches
	position       *position
}

func New(ctx context.Context, cfg *config.Config, runTypes []config.RunType) *Data {

	rt := make(map[config.RunType]struct{})
	for _, runType := range runTypes {
		rt[runType] = struct{}{}
	}

	return &Data{
		ctx: ctx,

		config: cfg,

		runTypes: rt,

		incomingDataCh:   make(chan types.Incoming, 2024),
		candleReceivedCh: make(chan struct{}, 1024),
		matchCh:          make(chan struct{}, 2024),
		positionCh:       make(chan struct{}, 512),

		candles:        newCandles(cfg.Data.CandlesSize),
		orderBookCache: newOrderBookCache(cfg.Data.OrderBookSize),
		matchesCache:   newMatches(cfg.Data.MatchesSize),
		position:       newPosition(),
	}
}

func (d *Data) GetCandles() *candles {
	return d.candles
}

func (d *Data) GetOrderBookCache() *orderBookCache {
	return d.orderBookCache
}

func (d *Data) GetMatches() *matches {
	return d.matchesCache
}

func (d *Data) GetPosition() *position {
	return d.position
}

func (d *Data) SendToIncomingCh(inc types.Incoming) {
	d.incomingDataCh <- inc
}

func (d *Data) IncomingCh() chan types.Incoming {
	return d.incomingDataCh
}

func (d *Data) sendToCandlesCh() {
	d.candleReceivedCh <- struct{}{}
}

func (d *Data) sendToMatchCh() {
	d.matchCh <- struct{}{}
}

func (d *Data) sendToPositionCh() {
	d.positionCh <- struct{}{}
}

func (d *Data) PositionCh() chan struct{} {
	return d.positionCh
}

func (d *Data) CandleCh() chan struct{} {
	return d.candleReceivedCh
}

func (d *Data) MatchCh() chan struct{} {
	return d.matchCh
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

				if _, ok := d.runTypes[config.EveryCandleRun]; ok {
					d.sendToCandlesCh()
				}

				d.candles.Set(candle)
			case types.IncomingOrderBook:
				ob := inc.OrderBook()

				d.orderBookCache.Set(ob)
			case types.IncomingMatch:
				p := inc.Match()

				if _, ok := d.runTypes[config.EveryMatchRun]; ok {
					d.sendToMatchCh()
				}

				d.matchesCache.Set(p)
			case types.IncomingOrder:
				o := inc.Order()

				if _, ok := d.runTypes[config.EveryPositionChangeRun]; ok {
					d.sendToPositionCh()
				}

				d.position.Set(o)
			default:
				panic("unknown incoming data")
			}
		}
	}
}

package data

import (
	"context"

	"github.com/tgmk/know-trade/util/sync"

	"github.com/tgmk/know-trade/config"
	"github.com/tgmk/know-trade/types"
)

type Data struct {
	ctx context.Context

	config *config.Config

	runTypes map[string]map[config.RunType]struct{}

	incomingDataCh   chan types.Incoming
	candleReceivedCh *sync.ChMap
	matchCh          *sync.ChMap
	positionCh       *sync.ChMap
	financeReportCh  chan struct{}

	candles        *candles
	orderBookCache *orderBookCache
	matchesCache   *matches
	position       *position
	finReports     *finReports
}

func New(ctx context.Context, cfg *config.Config, runTypes map[string]map[config.RunType]struct{}) *Data {

	return &Data{
		ctx: ctx,

		config: cfg,

		runTypes: runTypes,

		incomingDataCh:   make(chan types.Incoming, 2024),
		candleReceivedCh: sync.New(),
		matchCh:          sync.New(),
		positionCh:       sync.New(),
		financeReportCh:  make(chan struct{}, 256),

		candles:        newCandles(cfg.Data.CandlesSize),
		orderBookCache: newOrderBookCache(cfg.Data.OrderBookSize),
		matchesCache:   newMatches(cfg.Data.MatchesSize),
		position:       newPosition(),
		finReports:     newFinReports(),
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

func (d *Data) GetFinReports() *finReports {
	return d.finReports
}

func (d *Data) SendToIncomingCh(inc types.Incoming) {
	d.incomingDataCh <- inc
}

func (d *Data) IncomingCh() chan types.Incoming {
	return d.incomingDataCh
}

func (d *Data) sendToCandlesCh(instrument string) {
	ch := d.candleReceivedCh.Upsert(instrument)
	ch <- struct{}{}
}

func (d *Data) sendToMatchCh(instrument string) {
	ch := d.matchCh.Upsert(instrument)

	ch <- struct{}{}
}

func (d *Data) sendToPositionCh(instrument string) {
	ch := d.positionCh.Upsert(instrument)

	ch <- struct{}{}
}

func (d *Data) sendToFinReportsCh() {
	d.financeReportCh <- struct{}{}
}

func (d *Data) PositionCh(inst instrumentID) chan struct{} {
	return d.positionCh.Upsert(inst)
}

func (d *Data) CandleCh(inst instrumentID) chan struct{} {
	return d.candleReceivedCh.Upsert(inst)
}

func (d *Data) MatchCh(inst instrumentID) chan struct{} {
	return d.matchCh.Upsert(inst)
}

func (d *Data) FinReportsCh() chan struct{} {
	return d.financeReportCh
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

				d.candles.Set(candle)

				rt, ok := d.runTypes[candle.InstrumentID]
				if !ok {
					continue
				}

				if _, ok := rt[config.EveryCandleRun]; ok {
					d.sendToCandlesCh(candle.InstrumentID)
				}
			case types.IncomingOrderBook:
				ob := inc.OrderBook()

				d.orderBookCache.Set(ob)
			case types.IncomingMatch:
				p := inc.Match()

				d.matchesCache.Set(p)

				rt, ok := d.runTypes[p.InstrumentID]
				if !ok {
					continue
				}

				if _, ok := rt[config.EveryMatchRun]; ok {
					d.sendToMatchCh(p.InstrumentID)
				}
			case types.IncomingOrder:
				o := inc.Order()

				d.position.Set(o)

				go d.position.Update()

				rt, ok := d.runTypes[o.InstrumentID]
				if !ok {
					continue
				}

				if _, ok := rt[config.EveryPositionChangeRun]; ok {
					d.sendToPositionCh(o.InstrumentID)
				}
			case types.IncomingFinReport:
				fr := inc.FinReport()

				d.finReports.Set(fr)

				rt, ok := d.runTypes[fr.InstrumentID]
				if !ok {
					continue
				}

				if _, ok := rt[config.EveryFinReport]; ok {
					d.sendToFinReportsCh()
				}
			default:
				panic("unknown incoming data")
			}
		}
	}
}

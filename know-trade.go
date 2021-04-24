// Package knowtrade is designed to configure and run your trade strategies,
// without having to do routine operations
package knowtrade

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tgmk/know-trade/config"
	ctx "github.com/tgmk/know-trade/context"
	"github.com/tgmk/know-trade/data"
)

// Handler for implements your trade logic
type Handler func(ctx *ctx.Context, settings *RunSettings) error

// ErrHandler handles your trade strategies logic errors
type ErrHandler func(ctx *ctx.Context, err error) error

// Manager represents strategies runner
type Manager struct {
	ctx       *ctx.Context
	isStopped bool

	errCh chan error
	errH  ErrHandler

	howRun HowRun

	log logrus.FieldLogger
}

func New(
	ctx *ctx.Context,
	howRun HowRun,
	errH ErrHandler,
) *Manager {
	return &Manager{
		ctx: ctx,

		errCh: make(chan error),
		errH:  errH,

		howRun: howRun,

		log: logrus.New(),
	}
}

func (s *Manager) GetData() *data.Data {
	return s.ctx.GetData()
}

// Run runs your trade Manager logic
func (s *Manager) Run(errH ErrHandler) {
	go s.ctx.GetData().Process()

	if errH != nil {
		s.errH = errH
		go s.processErrors(s.errH)
	}

	for rt, runSettings := range s.howRun {
		runType := rt
		settings := runSettings
		switch runType {
		case config.TickerRun:
			go s.tickerRun(settings)
		case config.EveryCandleRun:
			go s.byCandleRun(settings)
		case config.EveryMatchRun:
			go s.byMatchRun(settings)
		case config.EveryPositionChangeRun:
			go s.byPositionRun(settings)
		case config.EveryFinReport:
			go s.byFinReportRun(settings)
		//case config.ByOthersRun:
		default:
			panic("unknown run type")
		}
	}
}

func (s *Manager) Stop() {
	if s.isStopped {
		return
	}

	s.ctx.CancelFunc()
	s.isStopped = true
}

func (s *Manager) tickerRun(settings *RunSettings) {
	ticker := time.NewTicker(settings.TickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			err := settings.Handler(s.ctx, settings)
			if err != nil {
				s.log.WithError(err).WithField("run", "ticker").Error("Manager execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *Manager) byCandleRun(settings *RunSettings) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().CandleCh(settings.InstrumentID):
			err := settings.Handler(s.ctx, settings)
			if err != nil {
				s.log.WithError(err).WithField("run", "every_candle").Error("Manager execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *Manager) byMatchRun(settings *RunSettings) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().MatchCh(settings.InstrumentID):
			err := settings.Handler(s.ctx, settings)
			if err != nil {
				s.log.WithError(err).WithField("run", "every_match").Error("Manager execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *Manager) byPositionRun(settings *RunSettings) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().PositionCh(settings.InstrumentID):
			err := settings.Handler(s.ctx, settings)
			if err != nil {
				s.log.WithError(err).WithField("run", "every_position_change").
					Error("Manager execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *Manager) byFinReportRun(settings *RunSettings) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().FinReportsCh():
			err := settings.Handler(s.ctx, settings)
			if err != nil {
				s.log.WithError(err).WithField("run", "every_fin_report_change").
					Error("Manager execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *Manager) byOthersRun(settings *RunSettings) {
	err := settings.Handler(s.ctx, settings)
	if err != nil {
		s.log.WithError(err).WithField("run", "other").Error("Manager execute failed")
		if s.errH != nil {
			s.errCh <- err
		}
	}
}

func (s *Manager) processErrors(errH ErrHandler) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.errCh:
			if err != nil {
				resultErr := errH(s.ctx, err)
				if resultErr != nil {
					s.log.WithError(err).Error("exit by error process")
					return
				}
			}
		}
	}
}

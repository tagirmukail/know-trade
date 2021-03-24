// knowtrade package is designed to configure and run your trade strategies,
// without having to do routine operations
package knowtrade

import (
	"context"
	"time"

	ctx "github.com/tgmk/know-trade/internal/context"

	"github.com/sirupsen/logrus"

	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/data"
)

// Handler for implements your trade logic
type Handler func(ctx context.Context, cfg *config.Config, d *data.Data) error

// ErrHandler handles your trade strategies logic errors
type ErrHandler func(ctx context.Context, cfg *config.Config, d *data.Data, err error) error

// strategy represents strategy runner
type strategy struct {
	ctx *ctx.Context

	errCh chan error
	errH  ErrHandler

	log logrus.FieldLogger
}

func New(
	ctx *ctx.Context,
) *strategy {
	return &strategy{

		errCh: make(chan error),

		log: logrus.New(),
	}
}

func (s *strategy) GetData() *data.Data {
	return s.ctx.GetData()
}

// Run runs your trade strategy logic
func (s *strategy) Run(ctx context.Context, h Handler, errH ErrHandler) {
	go s.ctx.GetData().Process()

	if errH != nil {
		s.errH = errH
		go s.processErrors(ctx, s.errH)
	}

	switch s.ctx.GetConfig().HowRun {
	case config.TickerRun:
		go s.tickerRun(ctx, h)
	case config.EveryCandleRun:
		go s.byCandleRun(ctx, h)
	case config.EveryPrintRun:
		go s.byPrintRun(ctx, h)
	case config.ByOthersRun:
	default:
		panic("unknown run type")
	}
}

func (s *strategy) tickerRun(ctx context.Context, h Handler) {
	ticker := time.NewTicker(s.ctx.GetConfig().TickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			err := h(ctx, s.ctx.GetConfig(), s.GetData())
			if err != nil {
				s.log.WithError(err).WithField("run", "ticker").Error("strategy execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *strategy) byCandleRun(ctx context.Context, h Handler) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().CandleCh():
			err := h(ctx, s.ctx.GetConfig(), s.GetData())
			if err != nil {
				s.log.WithError(err).WithField("run", "every_candle").Error("strategy execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *strategy) byPrintRun(ctx context.Context, h Handler) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().PrintCh():
			err := h(ctx, s.ctx.GetConfig(), s.GetData())
			if err != nil {
				s.log.WithError(err).WithField("run", "every_print").Error("strategy execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *strategy) byOthersRun(ctx context.Context, h Handler) {
	err := h(ctx, s.ctx.GetConfig(), s.GetData())
	if err != nil {
		s.log.WithError(err).WithField("run", "every_print").Error("strategy execute failed")
		if s.errH != nil {
			s.errCh <- err
		}
	}
}

func (s *strategy) processErrors(ctx context.Context, errH ErrHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-s.errCh:
			if err != nil {
				resultErr := errH(ctx, s.ctx.GetConfig(), s.GetData(), err)
				if resultErr != nil {
					s.log.WithError(err).WithField("run", s.ctx.GetConfig().HowRun).Error("exit by error process")
					return
				}
			}
		}
	}
}

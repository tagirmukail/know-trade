// knowtrade package is designed to configure and run your trade strategies,
// without having to do routine operations
package knowtrade

import (
	"time"

	ctx "github.com/tgmk/know-trade/internal/context"

	"github.com/sirupsen/logrus"

	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/data"
)

// Handler for implements your trade logic
type Handler func(ctx *ctx.Context) error

// ErrHandler handles your trade strategies logic errors
type ErrHandler func(ctx *ctx.Context, err error) error

// strategy represents strategy runner
type strategy struct {
	ctx *ctx.Context

	errCh chan error
	errH  ErrHandler

	log logrus.FieldLogger
}

func New(
	ctx *ctx.Context,
	errH ErrHandler,
) *strategy {
	return &strategy{
		ctx: ctx,

		errCh: make(chan error),
		errH:  errH,

		log: logrus.New(),
	}
}

func (s *strategy) GetData() *data.Data {
	return s.ctx.GetData()
}

// Run runs your trade strategy logic
func (s *strategy) Run(h Handler, errH ErrHandler) {
	go s.ctx.GetData().Process()

	if errH != nil {
		s.errH = errH
		go s.processErrors(s.errH)
	}

	switch s.ctx.GetConfig().HowRun {
	case config.TickerRun:
		go s.tickerRun(h)
	case config.EveryCandleRun:
		go s.byCandleRun(h)
	case config.EveryMatchRun:
		go s.byPrintRun(h)
	//case config.ByOthersRun:
	default:
		panic("unknown run type")
	}
}

func (s *strategy) tickerRun(h Handler) {
	ticker := time.NewTicker(s.ctx.GetConfig().TickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			err := h(s.ctx)
			if err != nil {
				s.log.WithError(err).WithField("run", "ticker").Error("strategy execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *strategy) byCandleRun(h Handler) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().CandleCh():
			err := h(s.ctx)
			if err != nil {
				s.log.WithError(err).WithField("run", "every_candle").Error("strategy execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *strategy) byPrintRun(h Handler) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.GetData().PrintCh():
			err := h(s.ctx)
			if err != nil {
				s.log.WithError(err).WithField("run", "every_print").Error("strategy execute failed")
				if s.errH != nil {
					s.errCh <- err
				}
			}
		}
	}
}

func (s *strategy) byOthersRun(h Handler) {
	err := h(s.ctx)
	if err != nil {
		s.log.WithError(err).WithField("run", "every_print").Error("strategy execute failed")
		if s.errH != nil {
			s.errCh <- err
		}
	}
}

func (s *strategy) processErrors(errH ErrHandler) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.errCh:
			if err != nil {
				resultErr := errH(s.ctx, err)
				if resultErr != nil {
					s.log.WithError(err).WithField("run", s.ctx.GetConfig().HowRun).Error("exit by error process")
					return
				}
			}
		}
	}
}

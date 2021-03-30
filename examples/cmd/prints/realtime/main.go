package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	appContext "github.com/tgmk/know-trade/internal/context"

	"github.com/Kucoin/kucoin-go-sdk"

	knowtrade "github.com/tgmk/know-trade"
	"github.com/tgmk/know-trade/examples/testcli"
	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/data"
	"github.com/tgmk/know-trade/internal/types"
)

func main() {
	var (
		feePercent   float64
		fee          float64
		startBalance float64
		pair         string
	)

	flag.Float64Var(&feePercent, "feePercent", 3, "fee in percent")
	flag.Float64Var(&fee, "fee", 0, "fixed fee")
	flag.Float64Var(&startBalance, "balance", 1000, "start balance")
	flag.StringVar(&pair, "pair", "BTC-USDT", "pair")

	flag.Parse()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	cli := testcli.New(ctx, fee, feePercent, startBalance)

	cfg := &config.Config{
		Run: config.Run{
			HowRun:       config.EveryMatchRun,
			InstrumentID: pair,
		},
		Data: config.Data{
			CandlesSize:   20,
			OrderBookSize: 20,
			MatchesSize:   120,
		},
	}

	aCtx := appContext.New(ctx, cfg)

	aCtx.SetExchangeClient(cli)

	s := knowtrade.New(aCtx, nil)

	d := aCtx.GetData()

	go readRealTimeFromExchange(ctx, pair, d)

	s.Run(strategyHandler, nil)

	<-done
	cancel()
}

func strategyHandler(ctx *appContext.Context) error {
	matches := ctx.GetData().GetMatches()

	lastMatch := matches.GetLast(ctx.GetConfig().InstrumentID)

	pair := lastMatch.InstrumentID

	switch {
	case lastMatch.Size > 0.1 && lastMatch.Side == "sell":
		o, err := ctx.GetExchangeClient().Limit(ctx, &types.LimitOrderRequest{
			Price:        lastMatch.Price,
			Size:         0.0001,
			InstrumentID: pair,
			Side:         "sell",
		})
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	case lastMatch.Size > 0.1 && lastMatch.Side == "buy":
		o, err := ctx.GetExchangeClient().Limit(ctx, &types.LimitOrderRequest{
			Price:        lastMatch.Price,
			Size:         0.0001,
			InstrumentID: pair,
			Side:         "buy",
		})
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	default:
		log.Printf("skip print: %#v", lastMatch)
	}

	cli, _ := ctx.GetExchangeClient().(*testcli.TestExchangeClient)

	r := cli.Result()

	log.Printf("balance: %v", r.Balance)
	log.Printf("earning: %v", r.Earning)

	return nil
}

func readRealTimeFromExchange(ctx context.Context, pair string, d *data.Data) {
	s := kucoin.NewApiService(
		kucoin.ApiKeyOption(""),
		kucoin.ApiSecretOption(""),
		kucoin.ApiPassPhraseOption(""),
		kucoin.ApiKeyVersionOption("2"),
	)

	rsp, err := s.WebSocketPublicToken()
	if err != nil {
		log.Fatal(err)
	}

	tk := &kucoin.WebSocketTokenModel{}
	err = rsp.ReadData(tk)
	if err != nil {
		log.Fatal(err)
	}

	c := s.NewWebSocketClient(tk)
	mc, ec, err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ch := kucoin.NewSubscribeMessage("/market/match:"+pair, false)

	err = c.Subscribe(ch)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			c.Stop()
			return
		case err = <-ec:
			c.Stop()
			if err != nil {
				log.Fatal(err)
			}
		case msg := <-mc:
			var match = &testcli.MatchMessage{}

			err = msg.ReadData(match)
			if err != nil {
				c.Stop()
				log.Fatal(err)
			}

			m := &types.Match{
				InstrumentID: match.GetSymbol(),
				Price:        match.GetPrice(),
				Size:         match.GetSize(),
				Side:         match.GetSide(),
				Time:         match.GetTime(),
				Taker:        match.TakerOrderID,
				Maker:        match.MakerOrderID,
			}

			d.SendToIncomingCh(m)

			log.Printf("print message sended: %#v", m)
		}
	}
}

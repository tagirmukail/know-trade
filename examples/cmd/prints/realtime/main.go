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
		symbol       string
	)

	flag.Float64Var(&feePercent, "feePercent", 3, "fee in percent")
	flag.Float64Var(&fee, "fee", 0, "fixed fee")
	flag.Float64Var(&startBalance, "balance", 1000, "start balance")
	flag.StringVar(&symbol, "symbol", "BTC-USDT", "symbol")

	flag.Parse()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	cli := testcli.New(ctx, fee, feePercent, startBalance)

	cfg := &config.Config{
		Run: config.Run{
			HowRun: config.EveryPrintRun,
		},
		Data: config.Data{
			CandlesSize:   20,
			OrderBookSize: 20,
			PrintsSize:    120,
		},
		ExchangeClient: cli,
	}

	aCtx := appContext.New(ctx, cfg)

	s := knowtrade.New(aCtx)

	d := aCtx.GetData()

	go readRealTimeFromExchange(ctx, symbol, d)

	s.Run(ctx, strategyHandler, nil)

	<-done
	cancel()
}

func strategyHandler(ctx context.Context, cfg *config.Config, d *data.Data) error {
	prints := d.GetPrints()

	lastPrint := prints.GetLast()

	symbol := lastPrint.Symbol

	switch {
	case lastPrint.Size > 0.1 && lastPrint.Side == "sell":
		o, err := cfg.ExchangeClient.Limit(ctx, symbol, "sell", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	case lastPrint.Size > 0.1 && lastPrint.Side == "buy":
		o, err := cfg.ExchangeClient.Limit(ctx, symbol, "buy", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	default:
		log.Printf("skip print: %#v", lastPrint)
	}

	return nil
}

func readRealTimeFromExchange(ctx context.Context, symbol string, d *data.Data) {
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

	ch := kucoin.NewSubscribeMessage("/market/match:"+symbol, false)

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

			p := &types.Print{
				Symbol: match.GetSymbol(),
				Price:  match.GetPrice(),
				Size:   match.GetSize(),
				Side:   match.GetSide(),
				Time:   match.GetTime(),
			}

			d.SendToIncomingCh(p)

			log.Printf("print message sended: %#v", p)
		}
	}
}

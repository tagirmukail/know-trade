package main

import (
	"context"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	knowtrade "github.com/tgmk/know-trade"
	"github.com/tgmk/know-trade/examples/testcli"
	"github.com/tgmk/know-trade/internal/config"
	"github.com/tgmk/know-trade/internal/types"

	appContext "github.com/tgmk/know-trade/internal/context"
)

func main() {
	var (
		feePercent   float64
		fee          float64
		startBalance float64
		filePath     string
	)

	flag.Float64Var(&feePercent, "feePercent", 3, "fee in percent")
	flag.Float64Var(&fee, "fee", 0, "fixed fee")
	flag.Float64Var(&startBalance, "balance", 1000, "start balance")
	flag.StringVar(&filePath, "fpath", "prints.csv", "download data to file by path")

	flag.Parse()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	cli := testcli.New(ctx, fee, feePercent, startBalance)

	cfg := &config.Config{
		Run: config.Run{
			InstrumentID: "BTC-USDT",
			HowRun:       config.EveryMatchRun,
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

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(f)
	reader.Comma = ';'

	s.Run(strategyHandler, nil)

	for {
		var record []string
		record, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		time.Sleep(500 * time.Millisecond)

		if len(record) != 7 {
			continue
		}

		match := &testcli.MatchMessage{
			Symbol:       record[0],
			Price:        record[1],
			Size:         record[2],
			Side:         record[3],
			Time:         record[4],
			MakerOrderID: record[5],
			TakerOrderID: record[6],
		}

		m := &types.Match{
			InstrumentID: match.GetSymbol(),
			Price:        match.GetPrice(),
			Size:         match.GetSize(),
			Side:         match.GetSide(),
			Time:         match.GetTime(),
			Maker:        match.MakerOrderID,
			Taker:        match.TakerOrderID,
		}

		d.SendToIncomingCh(m)
	}

	cancel()

	r := cli.Result()

	log.Printf("balance: %v", r.Balance)
	log.Printf("earning: %v", r.Earning)
}

func strategyHandler(ctx *appContext.Context) error {
	matches := ctx.GetData().GetMatches()

	lastMatch := matches.GetLast(ctx.GetConfig().InstrumentID)

	switch {
	case lastMatch.Size > 0.1 && lastMatch.Side == "sell":
		o, err := ctx.GetExchangeClient().Limit(ctx, &types.LimitOrderRequest{
			Price:        lastMatch.Price,
			Size:         0.0001,
			InstrumentID: ctx.GetConfig().InstrumentID,
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
			InstrumentID: ctx.GetConfig().InstrumentID,
			Side:         "buy",
		})
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	default:
		log.Printf("skip print: %#v", lastMatch)
	}

	return nil
}

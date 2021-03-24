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
	"github.com/tgmk/know-trade/internal/data"
	"github.com/tgmk/know-trade/internal/types"
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
			HowRun: config.EveryPrintRun,
		},
		Data: config.Data{
			CandlesSize:   20,
			OrderBookSize: 20,
			PrintsSize:    120,
		},
		ExchangeClient: cli,
	}

	s := knowtrade.New(ctx, cfg)

	d := s.GetData()

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(f)
	reader.Comma = ';'

	s.Run(ctx, strategyHandler, nil)

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

		if len(record) != 5 {
			continue
		}

		match := &testcli.MatchMessage{
			Symbol: record[0],
			Price:  record[1],
			Size:   record[2],
			Side:   record[3],
			Time:   record[4],
		}

		p := &types.Print{
			Symbol: match.GetSymbol(),
			Price:  match.GetPrice(),
			Size:   match.GetSize(),
			Side:   match.GetSide(),
			Time:   match.GetTime(),
		}

		d.SendToIncomingCh(p)
	}

	cancel()

	r := cli.Result()

	log.Printf("balance: %v", r.Balance)
	log.Printf("earning: %v", r.Earning)
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

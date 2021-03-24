# Go SDK for create your trade strategies


## Install

```
go get github.com/tgmk/know-trade
```

## Usage

### Real time application

This sample use kucoin exchange api.

```go
import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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

  // initialize your exchange client
	cli := testcli.New(ctx, fee, feePercent, startBalance)

  // configure your strategy
	cfg := &config.Config{
		Run: config.Run{
			HowRun: config.EveryPrintRun,
		},
		Data: config.Data{
			CandlesSize:   20,
			OrderBookSize: 20,
			PrintsSize:    120,
		},
		OrderExecutor: cli,
	}

	s := knowtrade.New(ctx, cfg)

	d := s.GetData()

  // read data from exchange
	go readRealTimeFromExchange(ctx, symbol, d)

  // run your strategy
	s.Run(ctx, strategyHandler, nil)

	<-done
	cancel()
}

// strategyHandler your strategy handler
func strategyHandler(ctx context.Context, cfg *config.Config, d *data.Data) error {
	prints := d.GetPrints()

	lastPrint := prints.GetLast()

	symbol := lastPrint.Symbol

	switch {
	case lastPrint.Size > 0.1 && lastPrint.Side == "sell":
		o, err := cfg.OrderExecutor.Limit(ctx, symbol, "sell", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	case lastPrint.Size > 0.1 && lastPrint.Side == "buy":
		o, err := cfg.OrderExecutor.Limit(ctx, symbol, "buy", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	default:
		log.Printf("skip print: %#v", lastPrint)
	}

	return nil
}

// readRealTimeFromExchange reads data from exchange and send to strategy executor
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
```

### Backtest

#### save data to file
```go

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"sync/atomic"

	"github.com/tgmk/know-trade/examples/testcli"

	"github.com/Kucoin/kucoin-go-sdk"
)

func main() {
	var (
		rows     int64
		filePath string
		symbol   string
	)

	flag.Int64Var(&rows, "rows", 2000, "download rows")
	flag.StringVar(&filePath, "fpath", "prints.csv", "download data to file by path")
	flag.StringVar(&symbol, "symbol", "BTC-USDT", "symbol")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	var csvWr *csv.Writer

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvWr = csv.NewWriter(f)
	csvWr.Comma = ';'
	defer csvWr.Flush()

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

	var canceled bool
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

			atomic.AddInt64(&rows, -1)

			record := []string{
				match.Symbol,
				match.Price,
				match.Size,
				match.Side,
				match.Time,
			}

			err = csvWr.Write(record)
			if err != nil {
				c.Stop()
				log.Fatal(err)
			}

			log.Printf("writed record to file: %v: %v", filePath, record)

			if atomic.LoadInt64(&rows) <= 0 {
				if canceled {
					return
				}

				cancel()
				canceled = true
				c.Stop()
				return
			}
		}
	}
}

```

#### Read data from file
```go
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
		OrderExecutor: cli,
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
		o, err := cfg.OrderExecutor.Limit(ctx, symbol, "sell", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	case lastPrint.Size > 0.1 && lastPrint.Side == "buy":
		o, err := cfg.OrderExecutor.Limit(ctx, symbol, "buy", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	default:
		log.Printf("skip print: %#v", lastPrint)
	}

	return nil
}

```

#### More examples see [examples](./examples/cmd)
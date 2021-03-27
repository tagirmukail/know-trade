# Go SDK for create your trade strategies(In developing)


## Install

```
go get github.com/tgmk/know-trade
```

## Usage
```go
config.Config{
		Run: config.Run{
            Symbol: "BTC-USDT" // you can run strategy for trade by this symbol or use stock identification
			HowRun: config.EveryPrintRun, // you can run your strategy by every incoming print(config.EveryPrintRun), 
			// every incoming candle(config.EveryPrintRun), by tick(config.TickerRun), install config.TickerInterval
			TickerInterval: time.Minute,
		},
		Data: config.Data{
			CandlesSize:   20, // install limit for candles in cache
			OrderBookSize: 20, // install limit for order book items in cache
			PrintsSize:    120, // install limit for prints(order matches) in cache
		},
	}
```

### Real time application

This sample use kucoin exchange api.

```go
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
	}

	aCtx := appContext.New(ctx, cfg)

	aCtx.SetExchangeClient(cli)

	s := knowtrade.New(aCtx, nil)

	d := aCtx.GetData()

	// read data from exchange
	go readRealTimeFromExchange(ctx, symbol, d)

	// run your strategy
	s.Run(strategyHandler, nil)

	<-done
	cancel()
}

// strategyHandler your strategy handler
func strategyHandler(ctx *appContext.Context) error {
	prints := ctx.GetData().GetPrints()

	lastPrint := prints.GetLast(ctx.GetConfig().Symbol)

	symbol := lastPrint.Symbol

	switch {
	case lastPrint.Size > 0.1 && lastPrint.Side == "sell":
		o, err := ctx.GetExchangeClient().Limit(ctx, symbol, "sell", lastPrint.Price, 0.0001)
		if err != nil {
			return err
		}

		log.Printf("executed: %#v", o)
	case lastPrint.Size > 0.1 && lastPrint.Side == "buy":
		o, err := ctx.GetExchangeClient().Limit(ctx, symbol, "buy", lastPrint.Price, 0.0001)
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

#### [save data to file](./examples/cmd/prints/download-prints/main.go)


#### [Run testing](./examples/cmd/prints/backtest/main.go)


#### More [examples](./examples/cmd)
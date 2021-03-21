package main

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

package main

import (
	"context"
	"log"
	"os"

	"github.com/tgmk/know-trade/finance/yahoo"
)

func main() {
	k := os.Getenv("key")
	if k == "" {
		log.Fatal("env `key` is empty")
	}

	finCli := yahoo.New("", k)

	resp, err := finCli.GetMarketQuotes(context.Background(), []string{"AAPL"}, yahoo.US)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range resp.QuoteResponse.Result {
		log.Printf("R: %#v", r)
	}

	anResp, err := finCli.GetStockAnalysis(context.Background(), "F", yahoo.US)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", anResp)
}

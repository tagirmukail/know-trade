package yahoo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// TODO add
// stock/v2/get-insider-transactions
// stock/v2/get-insights
// stock/v2/get-insider-roster
// stock/v2/get-holders
// stock/v2/get-balance-sheet
// stock/v2/get-cash-flow

const (
	apiKeyHeader         = "x-rapidapi-key"
	hostKeyHeader        = "x-rapidapi-host"
	contentTypeKeyHeader = "Content-Type"
)

type Region string

const (
	US Region = "US"
	BR Region = "BR"
	AU Region = "AU"
	CA Region = "CA"
	FR Region = "FR"
	DE Region = "DE"
	HK Region = "HK"
	IN Region = "IN"
	IT Region = "IT"
	ES Region = "ES"
	GB Region = "GB"
	SG Region = "SG"
)

const (
	requestsInSecondLimit = 5
	perSecondDuration     = time.Minute
)

type Client struct {
	client http.Client
	host   url.URL
	h      string
	key    string

	mx sync.Mutex

	perSecondStart time.Time
	perSecondCount int
}

func New(host, key string) *Client {
	h := url.URL{
		Scheme: "https",
		Host:   "apidojo-yahoo-finance-v1.p.rapidapi.com",
	}

	if host != "" {
		h.Host = host
	}

	client := http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 20 * time.Second,
			MaxIdleConns:    10,
		},
		Timeout: 30 * time.Second,
	}

	return &Client{
		host:   h,
		key:    key,
		client: client,
	}
}

// GetStockAnalysis gets data in Analysis section
func (c *Client) GetStockAnalysis(ctx context.Context, symbol string, region Region) (resp *StockResponse, err error) {
	resp = &StockResponse{}
	err = c.get(ctx, "/stock/v2/get-analysis", url.Values{
		"symbol": []string{symbol},
		"region": []string{string(region)},
	}, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetMarketQuotes the quotes by symbols
func (c *Client) GetMarketQuotes(ctx context.Context, symbols []string, region Region) (resp *QuoteResponse, err error) {
	resp = &QuoteResponse{}

	err = c.get(ctx, "/market/v2/get-quotes", url.Values{
		"symbols": symbols,
		"region":  []string{string(region)},
	}, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) get(ctx context.Context, path string, values url.Values, response interface{}) error {
	return c.do(ctx, request{
		method:    http.MethodGet,
		path:      path,
		urlValues: values,
		reqBody:   nil,
		response:  response,
	})
}

func (c *Client) checkPerSecond() error {
	now := time.Now()

	if c.perSecondStart.IsZero() {
		c.perSecondStart = now
	}

	if now.After(c.perSecondStart.Add(perSecondDuration)) {
		c.perSecondStart = now
		c.perSecondCount = 0
	}

	c.perSecondCount++

	if c.perSecondCount > requestsInSecondLimit {
		return errors.New("5 requests per second limit exceeded")
	}

	return nil
}

func (c *Client) checkMonthlyLimit(h http.Header) error {
	remainingVal := h.Get(XRatelimitRequestsRemainingHeader)
	if remainingVal == "" {
		return nil
	}

	remaining, err := strconv.Atoi(remainingVal)
	if err != nil {
		return err
	}

	if remaining <= 0 {
		requestsLimit := h.Get(XRatelimitRequestsLimitHeader)
		return fmt.Errorf("monthly limit %s exceeded", requestsLimit)
	}

	return nil
}

func (c *Client) do(ctx context.Context, req request) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	err := c.checkPerSecond()
	if err != nil {
		return err
	}

	uri := c.host
	uri.Path += req.path

	var body io.Reader
	if req.reqBody != nil {
		rData, err := json.Marshal(req.reqBody)
		if err != nil {
			return err
		}

		body = bytes.NewReader(rData)
	}

	if req.urlValues != nil {
		uri.RawQuery = req.urlValues.Encode()
	}

	r, err := http.NewRequestWithContext(ctx, req.method, uri.String(), body)
	if err != nil {
		return err
	}

	r.Header.Add(apiKeyHeader, c.key)
	r.Header.Add(hostKeyHeader, c.host.Host)

	if req.reqBody != nil {
		r.Header.Set(contentTypeKeyHeader, "application/json")
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}

	err = c.checkMonthlyLimit(resp.Header)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > http.StatusMultipleChoices {
		return fmt.Errorf("status is not ok: %v, desc: %v", resp.Status, string(data))
	}

	err = json.Unmarshal(data, req.response)
	if err != nil {
		return err
	}

	return nil
}

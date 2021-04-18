package yahoo

type QuoteResponse struct {
	QuoteResponse QResponse `json:"quoteResponse"`
}

type QResponse struct {
	Result []*Quote `json:"result"`
	Error  string   `json:"error"`
}

type Quote struct {
	Language                          string  `json:"language"`
	Region                            string  `json:"region"`
	QuoteType                         string  `json:"quoteType"`
	Triggerable                       bool    `json:"triggerable"`
	Currency                          string  `json:"currency"`
	FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
	PriceHint                         float64 `json:"priceHint"`
	TotalCash                         float64 `json:"totalCash"`
	FloatShares                       float64 `json:"floatShares"`
	Ebitda                            float64 `json:"ebitda"`
	ShortRatio                        float64 `json:"shortRatio"`
	PreMarketChange                   float64 `json:"preMarketChange"`
	PreMarketChangePercent            float64 `json:"preMarketChangePercent"`
	PreMarketTime                     int64   `json:"preMarketTime"`
	TargetPriceHigh                   float64 `json:"targetPriceHigh"`
	TargetPriceLow                    float64 `json:"targetPriceLow"`
	TargetPriceMean                   float64 `json:"targetPriceMean"`
	TargetPriceMedian                 float64 `json:"targetPriceMedian"`
	PreMarketPrice                    float64 `json:"preMarketPrice"`
	HeldPercentInsiders               float64 `json:"heldPercentInsiders"`
	HeldPercentInstitutions           float64 `json:"heldPercentInstitutions"`
	PostMarketChangePercent           float64 `json:"postMarketChangePercent"`
	PostMarketTime                    int64   `json:"postMarketTime"`
	PostMarketPrice                   float64 `json:"postMarketPrice"`
	PostMarketChange                  float64 `json:"postMarketChange"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketTime                 int64   `json:"regularMarketTime"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	RegularMarketDayRange             string  `json:"regularMarketDayRange"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketVolume               float64 `json:"regularMarketVolume"`
	SharesShort                       float64 `json:"sharesShort"`
	SharesShortPrevMonth              float64 `json:"sharesShortPrevMonth"`
	ShortPercentFloat                 float64 `json:"shortPercentFloat"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	Bid                               float64 `json:"bid"`
	Ask                               float64 `json:"ask"`
	BidSize                           float64 `json:"bidSize"`
	AskSize                           float64 `json:"askSize"`
	Exchange                          string  `json:"exchange"`
	Market                            string  `json:"market"`
	MessageBoardId                    string  `json:"messageBoardId"`
	FullExchangeName                  string  `json:"fullExchangeName"`
	ShortName                         string  `json:"shortName"`
	LongName                          string  `json:"longName"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	AverageDailyVolume3Month          float64 `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           float64 `json:"averageDailyVolume10Day"`
	Beta                              float64 `json:"beta"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	ExDividendDate                    int64   `json:"exDividendDate"`
	EarningsTimestamp                 int64   `json:"earningsTimestamp"`
	EarningsTimestampStart            int64   `json:"earningsTimestampStart"`
	EarningsTimestampEnd              int64   `json:"earningsTimestampEnd"`
	TrailingPE                        float64 `json:"trailingPE"`
	PegRatio                          float64 `json:"pegRatio"`
	DividendsPerShare                 float64 `json:"dividendsPerShare"`
	Revenue                           float64 `json:"revenue"`
	PriceToSales                      float64 `json:"priceToSales"`
	MarketState                       string  `json:"marketState"`
	EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
	EpsForward                        float64 `json:"epsForward"`
	EpsCurrentYear                    float64 `json:"epsCurrentYear"`
	EpsNextQuarter                    float64 `json:"epsNextQuarter"`
	PriceEpsCurrentYear               float64 `json:"priceEpsCurrentYear"`
	PriceEpsNextQuarter               float64 `json:"priceEpsNextQuarter"`
	SharesOutstanding                 float64 `json:"sharesOutstanding"`
	BookValue                         float64 `json:"bookValue"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	MarketCap                         float64 `json:"marketCap"`
	ForwardPE                         float64 `json:"forwardPE"`
	PriceToBook                       float64 `json:"priceToBook"`
	SourceInterval                    float64 `json:"sourceInterval"`
	ExchangeDataDelayedBy             int64   `json:"exchangeDataDelayedBy"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	Symbol                            string  `json:"symbol"`
}

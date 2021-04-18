package yahoo

type StockResponse struct {
	RecommendationTrend RecommendationTrend `json:"recommendationTrend"`
	FinancialsTemplate  FinancialsTemplate  `json:"financialsTemplate"`
	Price               Price               `json:"price"`
	EarningsHistory     struct {
		MaxAge  int64                 `json:"maxAge"`
		History []*EarningHistoryItem `json:"history"`
	} `json:"earningsHistory"`
	IndexTrend    IndexTrend    `json:"indexTrend"`
	FinancialData FinancialData `json:"financialData"`
	EarningsTrend struct {
		MaxAge int64        `json:"maxAge"`
		Trend  []*TrendItem `json:"trend"`
	} `json:"earningsTrend"`
	QuoteType               QuoteType     `json:"quoteType"`
	SummaryDetail           SummaryDetail `json:"summaryDetail"`
	Symbol                  string        `json:"symbol"`
	UpgradeDowngradeHistory struct {
		MaxAge  int64                          `json:"maxAge"`
		History []*UpgradeDowngradeHistoryItem `json:"history"`
	} `json:"upgradeDowngradeHistory"`
}

type UpgradeDowngradeHistoryItem struct {
	EpochGradeDate int64  `json:"epochGradeDate"`
	Firm           string `json:"firm"`
	ToGrade        string `json:"toGrade"`
	FromGrade      string `json:"fromGrade"`
	Action         string `json:"action"`
}

type SummaryDetail struct {
	PreviousClose                Raw         `json:"previousClose"`
	RegularMarketOpen            Raw         `json:"regularMarketOpen"`
	TwoHundredDayAverage         Raw         `json:"twoHundredDayAverage"`
	TrailingAnnualDividendYield  Raw         `json:"trailingAnnualDividendYield"`
	PayoutRatio                  Raw         `json:"payoutRatio"`
	RegularMarketDayHigh         Raw         `json:"regularMarketDayHigh"`
	NavPrice                     Raw         `json:"navPrice"`
	AverageDailyVolume10Day      RawWithLong `json:"averageDailyVolume10Day"`
	TotalAssets                  RawWithLong `json:"totalAssets"`
	RegularMarketPreviousClose   Raw         `json:"regularMarketPreviousClose"`
	FiftyDayAverage              Raw         `json:"fiftyDayAverage"`
	TrailingAnnualDividendRate   Raw         `json:"trailingAnnualDividendRate"`
	Open                         Raw         `json:"open"`
	AverageVolume10Days          RawWithLong `json:"averageVolume10Days"`
	ExpireDate                   RawWithLong `json:"expireDate"`
	Yield                        RawWithLong `json:"yield"`
	DividendRate                 Raw         `json:"dividendRate"`
	ExDividendDate               Raw         `json:"exDividendDate"`
	Beta                         Raw         `json:"beta"`
	CirculatingSupply            Raw         `json:"circulatingSupply"`
	StartDate                    Raw         `json:"startDate"`
	RegularMarketDayLow          Raw         `json:"regularMarketDayLow"`
	PriceHint                    RawWithLong `json:"priceHint"`
	Currency                     string      `json:"currency"`
	TrailingPE                   Raw         `json:"trailingPE"`
	RegularMarketVolume          RawWithLong `json:"regularMarketVolume"`
	MarketCap                    RawWithLong `json:"marketCap"`
	AverageVolume                RawWithLong `json:"averageVolume"`
	PriceToSalesTrailing12Months Raw         `json:"priceToSalesTrailing12Months"`
	DayLow                       Raw         `json:"dayLow"`
	Ask                          Raw         `json:"ask"`
	Volume                       RawWithLong `json:"volume"`
	FiftyTwoWeekHigh             Raw         `json:"fiftyTwoWeekHigh"`
	ForwardPE                    Raw         `json:"forwardPE"`
	FiveYearAvgDividendYield     Raw         `json:"fiveYearAvgDividendYield"`
	FiftyTwoWeekLow              Raw         `json:"fiftyTwoWeekLow"`
	Bid                          Raw         `json:"bid"`
	Tradeable                    bool        `json:"tradeable"`
	DividendYield                Raw         `json:"dividendYield"`
	BidSize                      RawWithLong `json:"bidSize"`
	DayHigh                      Raw         `json:"dayHigh"`
	MaxAge                       int64       `json:"maxAge"`
}

type QuoteType struct {
	Exchange                  string `json:"exchange"`
	ShortName                 string `json:"shortName"`
	LongName                  string `json:"longName"`
	ExchangeTimezoneName      string `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName string `json:"exchangeTimezoneShortName"`
	IsEsgPopulated            bool   `json:"isEsgPopulated"`
	GmtOffSetMilliseconds     string `json:"gmtOffSetMilliseconds"`
	QuoteType                 string `json:"quoteType"`
	Symbol                    string `json:"symbol"`
	MessageBoardId            string `json:"messageBoardId"`
	Market                    string `json:"market"`
}

type TrendItem struct {
	MaxAge           int64            `json:"maxAge"`
	Period           string           `json:"period"`
	EndDate          string           `json:"endDate"`
	Growth           Raw              `json:"growth"`
	EarningsEstimate EarningsEstimate `json:"earningsEstimate"`
	RevenueEstimate  RevenueEstimate  `json:"revenueEstimate"`
	EpsTrend         EpsTrend         `json:"epsTrend"`
	EpsRevisions     EpsRevisions     `json:"epsRevisions"`
}

type EpsRevisions struct {
	UpLast7Days    RawWithLong `json:"upLast7Days"`
	UpLast30Days   RawWithLong `json:"upLast30Days"`
	DownLast30Days RawWithLong `json:"downLast30Days"`
	DownLast90Days RawWithLong `json:"downLast90Days"`
}

type EpsTrend struct {
	Current       Raw `json:"current"`
	SevenDaysAgo  Raw `json:"7daysAgo"`
	ThirtyDaysAgo Raw `json:"30daysAgo"`
	SixtyDaysAgo  Raw `json:"60daysAgo"`
	NinetyDaysAgo Raw `json:"90daysAgo"`
}

type RevenueEstimate struct {
	Avg              Raw         `json:"avg"`
	Low              Raw         `json:"low"`
	High             Raw         `json:"high"`
	YearAgoRevenue   Raw         `json:"yearAgoRevenue"`
	NumberOfAnalysts RawWithLong `json:"numberOfAnalysts"`
	Growth           Raw         `json:"growth"`
}

type EarningsEstimate struct {
	Avg              Raw         `json:"avg"`
	Low              Raw         `json:"low"`
	High             Raw         `json:"high"`
	YearAgoEps       Raw         `json:"yearAgoEps"`
	NumberOfAnalysts RawWithLong `json:"numberOfAnalysts"`
	Growth           Raw         `json:"growth"`
}

type FinancialData struct {
	EbitdaMargins          Raw         `json:"ebitdaMargins"`
	ProfitMargins          Raw         `json:"profitMargins"`
	GrossMargins           Raw         `json:"grossMargins"`
	OperatingCashflow      RawWithLong `json:"operatingCashflow"`
	RevenueGrowth          Raw         `json:"revenueGrowth"`
	OperatingMargins       Raw         `json:"operatingMargins"`
	Ebitda                 RawWithLong `json:"ebitda"`
	TargetLowPrice         Raw         `json:"targetLowPrice"`
	RecommendationKey      string      `json:"recommendationKey"`
	GrossProfits           RawWithLong `json:"grossProfits"`
	FreeCashflow           RawWithLong `json:"freeCashflow"`
	TargetMedianPrice      Raw         `json:"targetMedianPrice"`
	CurrentPrice           Raw         `json:"currentPrice"`
	EarningsGrowth         Raw         `json:"earningsGrowth"`
	CurrentRatio           Raw         `json:"currentRatio"`
	ReturnOnAssets         Raw         `json:"returnOnAssets"`
	NumberOfAnalystOptions RawWithLong `json:"numberOfAnalystOptions"`
	TargetMeanPrice        Raw         `json:"targetMeanPrice"`
	DebtToEquity           Raw         `json:"debtToEquity"`
	ReturnOnEquity         RawWithLong `json:"returnOnEquity"`
	TargetHighPrice        Raw         `json:"targetHighPrice"`
	TotalCash              RawWithLong `json:"totalCash"`
	TotalDebt              RawWithLong `json:"totalDebt"`
	TotalRevenue           RawWithLong `json:"totalRevenue"`
	TotalCashPerShare      Raw         `json:"totalCashPerShare"`
	FinancialCurrency      string      `json:"financialCurrency"`
	MaxAge                 int64       `json:"maxAge"`
	RevenuePerShare        Raw         `json:"revenuePerShare"`
	QuickRatio             Raw         `json:"quickRatio"`
	RecommendationMean     Raw         `json:"recommendationMean"`
}

type IndexTrend struct {
	MaxAge    int64      `json:"maxAge"`
	Symbol    string     `json:"symbol"`
	PeRatio   Raw        `json:"peRatio"`
	PegRatio  Raw        `json:"pegRatio"`
	Estimates []Estimate `json:"estimates"`
}

type Estimate struct {
	Period string `json:"period"`
	Growth Raw    `json:"growth"`
}

type EarningHistoryItem struct {
	MaxAge          int64  `json:"maxAge"`
	EpsActual       Raw    `json:"epsActual"`
	EpsEstimate     Raw    `json:"epsEstimate"`
	EpsDifference   Raw    `json:"epsDifference"`
	SurprisePercent Raw    `json:"surprisePercent"`
	Quarter         Raw    `json:"quarter"`
	Period          string `json:"period"`
}

type RecommendationTrend struct {
	Trend  []*Trend `json:"trend"`
	MaxAge int64    `json:"maxAge"`
}

type Trend struct {
	Period     string  `json:"period"`
	StrongBuy  float64 `json:"strongBuy"`
	Buy        float64 `json:"buy"`
	Hold       float64 `json:"hold"`
	Sell       float64 `json:"sell"`
	StrongSell float64 `json:"strongSell"`
}

type FinancialsTemplate struct {
	Code   string `json:"code"`
	MaxAge int64  `json:"maxAge"`
}

type Price struct {
	QuoteSourceName            string      `json:"quoteSourceName"`
	RegularMarketOpen          Raw         `json:"regularMarketOpen"`
	AverageDailyVolume3Month   RawWithLong `json:"averageDailyVolume3Month"`
	Exchange                   string      `json:"exchange"`
	RegularMarketTime          int64       `json:"regularMarketTime"`
	RegularMarketDayHigh       Raw         `json:"regularMarketDayHigh"`
	ShortName                  string      `json:"shortName"`
	AverageDailyVolume10Day    RawWithLong `json:"averageDailyVolume10Day"`
	RegularMarketChange        Raw         `json:"regularMarketChange"`
	CurrencySymbol             string      `json:"currencySymbol"`
	RegularMarketPreviousClose Raw         `json:"regularMarketPreviousClose"`
	PostMarketTime             int64       `json:"postMarketTime"`
	ExchangeDataDelayedBy      float64     `json:"exchangeDataDelayedBy"`
	PostMarketChange           Raw         `json:"postMarketChange"`
	PostMarketPrice            Raw         `json:"postMarketPrice"`
	ExchangeName               string      `json:"exchangeName"`
	RegularMarketDayLow        Raw         `json:"regularMarketDayLow"`
	PriceHint                  RawWithLong `json:"priceHint"`
	Currency                   string      `json:"currency"`
	RegularMarketPrice         Raw         `json:"regularMarketPrice"`
	RegularMarketVolume        RawWithLong `json:"regularMarketVolume"`
	RegularMarketSource        string      `json:"regularMarketSource"`
	MarketState                string      `json:"marketState"`
	MarketCap                  RawWithLong `json:"marketCap"`
	QuoteType                  string      `json:"quoteType"`
	PostMarketSource           string      `json:"postMarketSource"`
	Symbol                     string      `json:"symbol"`
	PostMarketChangePercent    Raw         `json:"postMarketChangePercent"`
	PreMarketSource            string      `json:"preMarketSource"`
	MaxAge                     int64       `json:"maxAge"`
	RegularMarketChangePercent Raw         `json:"regularMarketChangePercent"`
}

type Raw struct {
	Raw float64 `json:"raw,omitempty"`
	Fmt string  `json:"fmt,omitempty"`
}

type RawWithLong struct {
	Raw     float64 `json:"raw,omitempty"`
	Fmt     string  `json:"fmt,omitempty"`
	LongFmt string  `json:"longFmt,omitempty"`
}

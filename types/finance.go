package types

// https://journal.tinkoff.ru/analiz-emitenta/
// https://journal.tinkoff.ru/multilplicator/

type FinReport struct {
	*baseIncoming
	InstrumentID    string
	Period          string  // quarter, year ...
	MarketCap       float64 // market capitalization
	P               float64 // company price
	E               float64 // company earning
	PToE            float64 // price to earnings
	B               float64 // book ratio
	PToB            float64 // price to book ratio
	PToS            float64 // price to sales - it is the ratio of the market price of a share to the revenue per share
	Debt            float64
	Equity          float64
	DebtToEq        float64 // debt to equity ratio
	DSI             float64 // dividend Stability Index
	Payout          float64 // share of profits that the company pays for dividends
	QuickRatio      float64 // this is the ratio of highly liquid assets minus stocks to short-term liabilities
	CurrentRatio    float64 // unlike Quick Ratio, it takes into account hard-to-sell warehouse stocks
	NetProfitMargin float64 // ratio of net profit to revenue
	ROE             float64 // the indicator characterizes the efficiency of using shareholders' funds on an annualized basis
	PEG             float64 // P/E ratio to projected earnings growth
	ForwardPToE     float64 // reflects the expectation of the company's profit growth
	EPS             float64 // net earnings per share
	EV              float64 // enterprise value - this is the fair value of the company. EV = MarketCap + Debt - Equity
	EBITDA          float64 // this is the company's profit before interest, taxes and depreciation
	EVToEBITDA      float64 // this is the market value of the unit of profit
	DebtToEBITDA    float64 // reflects the number of years it takes a company to pay off all debts with its profit
}

func (q *FinReport) Type() IncomingType {
	return IncomingFinReport
}

func (q *FinReport) FinReport() *FinReport {
	return q
}

package yahoo

import "time"

// IntervalAPI represents a set of predefined time intervals used for specifying durations in API requests.
type IntervalAPI string

// Interval1m represents a 1-minute interval for the IntervalAPI type.
// Interval5m represents a 5-minute interval for the IntervalAPI type.
// Interval1d represents a 1-day interval for the IntervalAPI type.
// Interval1wk represents a 1-week interval for the IntervalAPI type.
// Interval1mo represents a 1-month interval for the IntervalAPI type.
const (
	Interval1m  IntervalAPI = "1m"  // one minute
	Interval5m  IntervalAPI = "5m"  // 5 minutes
	Interval1d  IntervalAPI = "1d"  // one day
	Interval1wk IntervalAPI = "1wk" // one week
	Interval1mo IntervalAPI = "1mo" // one month
)

// PeriodAPI represents a string type used to define specific periods in financial or date-related contexts.
type PeriodAPI string

// Period1d represents a period of 1 day.
// Period5d represents a period of 5 days.
// Period1mo represents a period of 1 month.
// Period3mo represents a period of 3 months.
// Period6mo represents a period of 6 months.
// Period1y represents a period of 1 year.
// Period2y represents a period of 2 years.
// Period5y represents a period of 5 years.
// Period10y represents a period of 10 years.
// PeriodYTD represents the year-to-date period.
// PeriodMax represents the maximum available period.
const (
	Period1d  PeriodAPI = "1d"  // 1 day
	Period5d  PeriodAPI = "5d"  // 5 days
	Period1mo PeriodAPI = "1mo" // 1 month
	Period3mo PeriodAPI = "3mo" // 3 months
	Period6mo PeriodAPI = "6mo" // 6 months
	Period1y  PeriodAPI = "1y"  // 1 year
	Period2y  PeriodAPI = "2y"  // 2 years
	Period5y  PeriodAPI = "5y"  // 5 years
	Period10y PeriodAPI = "10y" // 10 years
	PeriodYTD PeriodAPI = "ytd" // year-to-date
	PeriodMax PeriodAPI = "max" // maximum available
)

// YahooFinanceResponse represents the response from Yahoo Finance API
type YahooFinanceResponse struct {
	Chart ChartResponse `json:"chart"`
}

// ChartResponse represents the chart data in the Yahoo Finance response
type ChartResponse struct {
	Result []Result `json:"result"`
	Error  *Error   `json:"error"`
}

// Error represents an error in the Yahoo Finance response
type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Result represents a result in the Yahoo Finance response
type Result struct {
	Meta       Meta       `json:"meta"`
	Timestamp  []int64    `json:"timestamp"`
	Indicators Indicators `json:"indicators"`
}

// Meta represents metadata in the Yahoo Finance response
type Meta struct {
	Currency             string        `json:"currency"`
	Symbol               string        `json:"symbol"`
	Name                 string        `json:"longName"`
	ExchangeName         string        `json:"exchangeName"`
	InstrumentType       string        `json:"instrumentType"`
	FirstTradeDate       int64         `json:"firstTradeDate"`
	RegularMarketTime    int64         `json:"regularMarketTime"`
	GMTOffset            int           `json:"gmtoffset"`
	Timezone             string        `json:"timezone"`
	ExchangeTimezoneName string        `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64       `json:"regularMarketPrice"`
	ChartPreviousClose   float64       `json:"chartPreviousClose"`
	PriceHint            int           `json:"priceHint"`
	CurrentTradingPeriod TradingPeriod `json:"currentTradingPeriod"`
	DataGranularity      string        `json:"dataGranularity"`
	Range                string        `json:"range"`
	ValidRanges          []string      `json:"validRanges"`
}

// TradingPeriod represents a trading period in the Yahoo Finance response
type TradingPeriod struct {
	Pre     Period `json:"pre"`
	Regular Period `json:"regular"`
	Post    Period `json:"post"`
}

// Period represents a time period in the Yahoo Finance response
type Period struct {
	Timezone  string `json:"timezone"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	GMTOffset int    `json:"gmtoffset"`
}

// Indicators represents the indicators in the Yahoo Finance response
type Indicators struct {
	Quote    []Quote    `json:"quote"`
	Adjclose []Adjclose `json:"adjclose,omitempty"`
}

// Quote represents quote data in the Yahoo Finance response
type Quote struct {
	High   []float64 `json:"high"`
	Open   []float64 `json:"open"`
	Low    []float64 `json:"low"`
	Close  []float64 `json:"close"`
	Volume []float64 `json:"volume"`
}

// Adjclose represents adjusted close data in the Yahoo Finance response
type Adjclose struct {
	Adjclose []float64 `json:"adjclose"`
}

type StockPrice struct {
	Time     time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   int
}

type MarketData struct {
	Symbol   string
	Name     string
	Exchange string
	Prices   []StockPrice
}

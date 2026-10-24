package models

type FundamentalMetricsRequest struct {
	SecurityID string `json:"securityId"`
}


type FundamentalMetricsResponse struct {
	CompanyName string `json:"companyName"`
	Industry    string `json:"industry"`
	Description string `json:"description"`
}


type TopInstrumentsRequest struct {
	ExchangeSegment string `json:"exchangeSegment,omitempty"`
}


type TopInstrumentItem struct {
	TradingSymbol string  `json:"tradingSymbol"`
	SecurityID    string  `json:"securityId"`
	Change        float64 `json:"change"`
}


type TopInstrumentsResponse struct {
	Gainers []TopInstrumentItem `json:"gainers"`
	Losers  []TopInstrumentItem `json:"losers"`
}


type NewsHeadlineRequest struct {
	SecurityID string `json:"securityId,omitempty"`
}


type NewsItem struct {
	Title       string `json:"title"`
	Source      string `json:"source"`
	PublishedAt string `json:"publishedAt"`
}


type NewsHeadlineResponse struct {
	Headlines []NewsItem `json:"headlines"`
}


type TechnicalMetricsRequest struct {
	SecurityID      string `json:"securityId,omitempty"`
	ExchangeSegment string `json:"exchangeSegment,omitempty"`
}


type TechnicalMetricsResponse struct {
	SMA50  float64 `json:"sma50"`
	EMA20  float64 `json:"ema20"`
	RSI14  float64 `json:"rsi14"`
	PivotP float64 `json:"pivotP"`
}


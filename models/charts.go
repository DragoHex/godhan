package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type HistoricalChartsRequest struct {
	SecurityID      string `json:"securityId"`
	ExchangeSegment string `json:"exchangeSegment"`
	Instrument      string `json:"instrument"`
	ExpiryCode      int32  `json:"expiryCode,omitempty"`
	FromDate        string `json:"fromDate"`
	ToDate          string `json:"toDate"`
}


type IntradayChartsRequest struct {
	SecurityID      string `json:"securityId"`
	ExchangeSegment string `json:"exchangeSegment"`
	Instrument      string `json:"instrument"`
	Interval        string `json:"interval"` // 1, 5, 15, 25, 60
	FromDate        string `json:"fromDate"`
	ToDate          string `json:"toDate"`
}


type OptionChartRequest struct {
	SecurityID      string `json:"securityId"`
	ExchangeSegment string `json:"exchangeSegment"`
	Instrument      string `json:"instrument"`
	ExpiryDate      string `json:"expiryDate"`
	FromDate        string `json:"fromDate"`
	ToDate          string `json:"toDate"`
}


type ChartsResponse struct {
	Open         []float64 `json:"open"`
	High         []float64 `json:"high"`
	Low          []float64 `json:"low"`
	Close        []float64 `json:"close"`
	Volume       []int64   `json:"volume"`
	Timestamp    []int64   `json:"timestamp"`
	OpenInterest []int64   `json:"open_interest"`
}


type OptionChartResponse struct {
	Open      []float64 `json:"open"`
	High      []float64 `json:"high"`
	Low       []float64 `json:"low"`
	Close     []float64 `json:"close"`
	Volume    []int64   `json:"volume"`
	Timestamp []int64   `json:"timestamp"`
}


type Candle struct {
	Timestamp    time.Time
	Open         decimal.Decimal
	High         decimal.Decimal
	Low          decimal.Decimal
	Close        decimal.Decimal
	Volume       int64
	OpenInterest int64
}



// ToCandles converts parallel arrays into a slice of typed Candle structs with decimal precision
func (cr *ChartsResponse) ToCandles() []Candle {
	if cr == nil || len(cr.Open) == 0 {
		return nil
	}
	n := len(cr.Open)
	candles := make([]Candle, n)
	for i := 0; i < n; i++ {
		var ts time.Time
		if i < len(cr.Timestamp) {
			ts = time.Unix(cr.Timestamp[i], 0)
		}
		var vol int64
		if i < len(cr.Volume) {
			vol = cr.Volume[i]
		}
		var oi int64
		if i < len(cr.OpenInterest) {
			oi = cr.OpenInterest[i]
		}

		candles[i] = Candle{
			Timestamp:    ts,
			Open:         decimal.NewFromFloat(cr.Open[i]),
			High:         decimal.NewFromFloat(cr.High[i]),
			Low:          decimal.NewFromFloat(cr.Low[i]),
			Close:        decimal.NewFromFloat(cr.Close[i]),
			Volume:       vol,
			OpenInterest: oi,
		}
	}
	return candles
}

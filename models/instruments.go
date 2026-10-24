package models

import (
	"github.com/shopspring/decimal"
)

type Instrument struct {
	Exchange       string          `json:"exchange"`
	Segment        string          `json:"segment"`
	SecurityID     string          `json:"securityId"`
	InstrumentName string          `json:"instrumentName"`
	ExpiryDate     string          `json:"expiryDate"`
	TradingSymbol  string          `json:"tradingSymbol"`
	LotSize        int32           `json:"lotSize"`
	CustomSymbol   string          `json:"customSymbol"`
	StrikePrice    decimal.Decimal `json:"strikePrice"`
}


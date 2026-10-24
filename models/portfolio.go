package models

import (
	"github.com/shopspring/decimal"
)

type Holding struct {
	DhanClientID  string          `json:"dhanClientId"`
	TradingSymbol string          `json:"tradingSymbol"`
	SecurityID    string          `json:"securityId"`
	ISIN          string          `json:"isin"`
	TotalQty      int32           `json:"totalQty"`
	DpQty         int32           `json:"dpQty"`
	T1Qty         int32           `json:"t1Qty"`
	AvailableQty  int32           `json:"availableQty"`
	CollateralQty int32           `json:"collateralQty"`
	AvgCostPrice  decimal.Decimal `json:"avgCostPrice"`
}


type Position struct {
	DhanClientID     string          `json:"dhanClientId"`
	PositionType     string          `json:"positionType"` // LONG, SHORT, CLOSED
	TradingSymbol    string          `json:"tradingSymbol"`
	SecurityID       string          `json:"securityId"`
	ExchangeSegment  string          `json:"exchangeSegment"`
	ProductType      string          `json:"productType"`
	BuyAvg           decimal.Decimal `json:"buyAvg"`
	BuyQty           int32           `json:"buyQty"`
	SellAvg          decimal.Decimal `json:"sellAvg"`
	SellQty          int32           `json:"sellQty"`
	NetQty           int32           `json:"netQty"`
	RealizedProfit   decimal.Decimal `json:"realizedProfit"`
	UnrealizedProfit decimal.Decimal `json:"unrealizedProfit"`
}


type PositionConversionRequest struct {
	DhanClientID    string `json:"dhanClientId,omitempty"`
	FromProductType string `json:"fromProductType"`
	ToProductType   string `json:"toProductType"`
	ExchangeSegment string `json:"exchangeSegment"`
	PositionType    string `json:"positionType"`
	SecurityID      string `json:"securityId"`
	ConvertQty      int32  `json:"convertQty"`
}


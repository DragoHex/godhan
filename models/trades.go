package models

import (
	"github.com/shopspring/decimal"
)

type TradeResponse struct {
	DhanClientID    string          `json:"dhanClientId"`
	OrderID         string          `json:"orderId"`
	ExchangeOrderID string          `json:"exchangeOrderId"`
	ExchangeTradeID string          `json:"exchangeTradeId"`
	TransactionType string          `json:"transactionType"` // BUY, SELL
	ExchangeSegment string          `json:"exchangeSegment"` // NSE_EQ, NSE_FNO, BSE_EQ, BSE_FNO, MCX_COMM
	ProductType     string          `json:"productType"`     // CNC, INTRADAY, MARGIN, MTF, CO, BO
	OrderType       string          `json:"orderType"`
	TradingSymbol   string          `json:"tradingSymbol"`
	CustomSymbol    string          `json:"customSymbol"`
	SecurityID      string          `json:"securityId"`
	TradedQuantity  int32           `json:"tradedQuantity"`
	TradedPrice     decimal.Decimal `json:"tradedPrice"`
	CreateTime      string          `json:"createTime"`
	UpdateTime      string          `json:"updateTime"`
	ExchangeTime    string          `json:"exchangeTime"`
	DrvExpiryDate   string          `json:"drvExpiryDate"`
	DrvOptionType   string          `json:"drvOptionType"`
	DrvStrikePrice  decimal.Decimal `json:"drvStrikePrice"`
}


// KiteSDK Parity Aliases
type Trade = TradeResponse


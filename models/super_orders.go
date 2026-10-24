package models

import (
	"github.com/shopspring/decimal"
)

type SuperOrderRequest struct {
	DhanClientID    string          `json:"dhanClientId,omitempty"`
	CorrelationID   string          `json:"correlationId,omitempty"`
	TransactionType string          `json:"transactionType"` // BUY, SELL
	ExchangeSegment string          `json:"exchangeSegment"` // NSE_EQ, NSE_FNO, BSE_EQ, BSE_FNO, MCX_COMM
	ProductType     string          `json:"productType"`     // INTRADAY, MARGIN, MTF, CNC
	OrderType       string          `json:"orderType"`       // LIMIT, MARKET
	SecurityID      string          `json:"securityId"`
	Quantity        int32           `json:"quantity"`
	Price           decimal.Decimal `json:"price,omitempty"`
	TargetPrice     decimal.Decimal `json:"targetPrice"`
	StopLossPrice   decimal.Decimal `json:"stopLossPrice"`
	TrailingJump    decimal.Decimal `json:"trailingJump,omitempty"`
}


type SuperModifyRequest struct {
	DhanClientID  string          `json:"dhanClientId,omitempty"`
	OrderID       string          `json:"orderId,omitempty"`
	OrderType     string          `json:"orderType"`
	LegName       string          `json:"legName"` // ENTRY_LEG, STOP_LOSS_LEG, TARGET_LEG
	Quantity      int32           `json:"quantity"`
	Price         decimal.Decimal `json:"price,omitempty"`
	TargetPrice   decimal.Decimal `json:"targetPrice,omitempty"`
	StopLossPrice decimal.Decimal `json:"stopLossPrice,omitempty"`
	TrailingJump  decimal.Decimal `json:"trailingJump,omitempty"`
}


type SuperOrderResponse struct {
	DhanClientID    string          `json:"dhanClientId"`
	OrderID         string          `json:"orderId"`
	CorrelationID   string          `json:"correlationId"`
	OrderStatus     string          `json:"orderStatus"`
	TransactionType string          `json:"transactionType"`
	ExchangeSegment string          `json:"exchangeSegment"`
	ProductType     string          `json:"productType"`
	OrderType       string          `json:"orderType"`
	SecurityID      string          `json:"securityId"`
	Quantity        int32           `json:"quantity"`
	Price           decimal.Decimal `json:"price"`
	TargetPrice     decimal.Decimal `json:"targetPrice"`
	StopLossPrice   decimal.Decimal `json:"stopLossPrice"`
	TrailingJump    decimal.Decimal `json:"trailingJump"`
	FilledQty       int32           `json:"filledQty"`
}


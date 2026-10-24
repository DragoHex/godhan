package models

import (
	"github.com/shopspring/decimal"
)

type PostbackPayload struct {
	DhanClientID        string          `json:"dhanClientId"`
	OrderID             string          `json:"orderId"`
	CorrelationID       string          `json:"correlationId"`
	OrderStatus         string          `json:"orderStatus"`
	TransactionType     string          `json:"transactionType"`
	ExchangeSegment     string          `json:"exchangeSegment"`
	ProductType         string          `json:"productType"`
	OrderType           string          `json:"orderType"`
	Validity            string          `json:"validity"`
	TradingSymbol       string          `json:"tradingSymbol"`
	SecurityID          string          `json:"securityId"`
	Quantity            int32           `json:"quantity"`
	DisclosedQuantity   int32           `json:"disclosedQuantity"`
	Price               decimal.Decimal `json:"price"`
	TriggerPrice        decimal.Decimal `json:"triggerPrice"`
	AfterMarketOrder    bool            `json:"afterMarketOrder"`
	RemainingQuantity   int32           `json:"remainingQuantity"`
	FilledQty           int32           `json:"filledQty"`
	AverageTradedPrice  decimal.Decimal `json:"averageTradedPrice"`
	CreateTime          string          `json:"createTime"`
	UpdateTime          string          `json:"updateTime"`
	ExchangeTime        string          `json:"exchangeTime"`
	OmsErrorCode        string          `json:"omsErrorCode"`
	OmsErrorDescription string          `json:"omsErrorDescription"`
}


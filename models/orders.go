package models

import (
	"github.com/shopspring/decimal"
)

type OrderRequest struct {
	DhanClientID      string          `json:"dhanClientId,omitempty"`
	CorrelationID     string          `json:"correlationId,omitempty"`
	TransactionType   string          `json:"transactionType"`   // BUY, SELL
	ExchangeSegment   string          `json:"exchangeSegment"`   // NSE_EQ, NSE_FNO, BSE_EQ, BSE_FNO, MCX_COMM
	ProductType       string          `json:"productType"`       // CNC, INTRADAY, MARGIN, MTF, CO, BO
	OrderType         string          `json:"orderType"`         // LIMIT, MARKET, STOP_LOSS, STOP_LOSS_MARKET
	Validity          string          `json:"validity"`          // DAY, IOC
	TradingSymbol     string          `json:"tradingSymbol,omitempty"`
	SecurityID        string          `json:"securityId"`
	Quantity          int32           `json:"quantity"`
	DisclosedQuantity int32           `json:"disclosedQuantity,omitempty"`
	Price             decimal.Decimal `json:"price,omitempty"`
	TriggerPrice      decimal.Decimal `json:"triggerPrice,omitempty"`
	AfterMarketOrder  bool            `json:"afterMarketOrder,omitempty"`
	AMO               bool            `json:"amo,omitempty"`
	BoProfitValue     decimal.Decimal `json:"boProfitValue,omitempty"`
	BoStopLossValue   decimal.Decimal `json:"boStopLossValue,omitempty"`
}


type OrderModifyRequest struct {
	DhanClientID      string          `json:"dhanClientId,omitempty"`
	OrderID           string          `json:"orderId,omitempty"`
	OrderType         string          `json:"orderType"`
	LegName           string          `json:"legName,omitempty"`
	Quantity          int32           `json:"quantity"`
	Price             decimal.Decimal `json:"price,omitempty"`
	DisclosedQuantity int32           `json:"disclosedQuantity,omitempty"`
	TriggerPrice      decimal.Decimal `json:"triggerPrice,omitempty"`
	Validity          string          `json:"validity"`
}


type OrderStatusResponse struct {
	OrderID     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"` // TRANSIT, PENDING, REJECTED, CANCELLED, TRADED, EXPIRED
}


type OrderResponse struct {
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


// KiteSDK Parity Aliases
type OrderParams = OrderRequest
type OrderModifyParams = OrderModifyRequest
type Order = OrderResponse


package models

import (
	"github.com/shopspring/decimal"
)

type GTTOrderModel struct {
	DhanClientID      string          `json:"dhanClientId,omitempty"`
	CorrelationID     string          `json:"correlationId,omitempty"`
	OrderFlag         string          `json:"orderFlag"`         // SINGLE, OCO
	TransactionType   string          `json:"transactionType"`   // BUY, SELL
	ExchangeSegment   string          `json:"exchangeSegment"`   // NSE_EQ, NSE_FNO, BSE_EQ, BSE_FNO
	ProductType       string          `json:"productType"`       // CNC, MTF, MARGIN
	OrderType         string          `json:"orderType"`         // LIMIT, MARKET
	Validity          string          `json:"validity,omitempty"`
	SecurityID        string          `json:"securityId"`
	Quantity          int32           `json:"quantity"`
	DisclosedQuantity int32           `json:"disclosedQuantity,omitempty"`
	Price             decimal.Decimal `json:"price"`
	TriggerPrice      decimal.Decimal `json:"triggerPrice"`
	Price1            decimal.Decimal `json:"price1,omitempty"`        // OCO Target
	TriggerPrice1     decimal.Decimal `json:"triggerPrice1,omitempty"` // OCO Target Trigger
	Quantity1         int32           `json:"quantity1,omitempty"`     // OCO Target Quantity
}


type GttModifyRequest struct {
	DhanClientID  string          `json:"dhanClientId,omitempty"`
	OrderID       string          `json:"orderId,omitempty"`
	OrderFlag     string          `json:"orderFlag"`
	OrderType     string          `json:"orderType"`
	Quantity      int32           `json:"quantity"`
	Price         decimal.Decimal `json:"price"`
	TriggerPrice  decimal.Decimal `json:"triggerPrice"`
	Price1        decimal.Decimal `json:"price1,omitempty"`
	TriggerPrice1 decimal.Decimal `json:"triggerPrice1,omitempty"`
	Quantity1     int32           `json:"quantity1,omitempty"`
}


type GttOrderStatusResponse struct {
	OrderID     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"` // PENDING, REJECTED, CANCELLED, EXPIRED
}


type ForeverOrderResponse struct {
	DhanClientID    string          `json:"dhanClientId"`
	OrderID         string          `json:"orderId"`
	CorrelationID   string          `json:"correlationId"`
	OrderFlag       string          `json:"orderFlag"`
	OrderStatus     string          `json:"orderStatus"`
	TransactionType string          `json:"transactionType"`
	ExchangeSegment string          `json:"exchangeSegment"`
	ProductType     string          `json:"productType"`
	OrderType       string          `json:"orderType"`
	SecurityID      string          `json:"securityId"`
	Quantity        int32           `json:"quantity"`
	Price           decimal.Decimal `json:"price"`
	TriggerPrice    decimal.Decimal `json:"triggerPrice"`
}


type AlertCondition struct {
	ExchangeSegment string          `json:"exchangeSegment,omitempty"`
	SecurityID      string          `json:"securityId,omitempty"`
	Indicator       string          `json:"indicator,omitempty"`
	ConditionType   string          `json:"conditionType,omitempty"`
	ThresholdPrice  decimal.Decimal `json:"thresholdPrice,omitempty"`
}


type AlertOrderRequest struct {
	DhanClientID string          `json:"dhanClientId,omitempty"`
	Condition    *AlertCondition `json:"condition,omitempty"`
	Orders       []OrderRequest  `json:"orders,omitempty"`
}


type AlertModifyRequest struct {
	Condition *AlertCondition `json:"condition,omitempty"`
	Orders    []OrderRequest  `json:"orders,omitempty"`
}


type AlertOrderResponse struct {
	AlertID     string `json:"alertId"`
	AlertStatus string `json:"alertStatus"` // ACTIVE, TRIGGERED, EXPIRED, CANCELLED
}


type GetAlertResponse struct {
	AlertID     string          `json:"alertId"`
	AlertStatus string          `json:"alertStatus"`
	Condition   *AlertCondition `json:"condition,omitempty"`
	Orders      []OrderRequest  `json:"orders,omitempty"`
}


type MultiOrderRequest struct {
	DhanClientID string         `json:"dhanClientId,omitempty"`
	Orders       []OrderRequest `json:"orders"`
}


type MultiOrderResponse struct {
	MultiOrderID string `json:"multiOrderId"`
	Status       string `json:"status"`
}


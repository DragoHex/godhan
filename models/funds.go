package models

import (
	"github.com/shopspring/decimal"
)

type FundLimitResponse struct {
	DhanClientID        string          `json:"dhanClientId"`
	AvailableBalance    decimal.Decimal `json:"availabelBalance"`
	SodLimit            decimal.Decimal `json:"sodLimit"`
	CollateralAmount    decimal.Decimal `json:"collateralAmount"`
	ReceiveableAmount   decimal.Decimal `json:"receiveableAmount"`
	UtilizedAmount      decimal.Decimal `json:"utilizedAmount"`
	BlockedPayoutAmount decimal.Decimal `json:"blockedPayoutAmount"`
	WithdrawableBalance decimal.Decimal `json:"withdrawableBalance"`
}


type KnowYourMarginReq struct {
	DhanClientID    string          `json:"dhanClientId,omitempty"`
	ExchangeSegment string          `json:"exchangeSegment"`
	TransactionType string          `json:"transactionType"`
	Quantity        int32           `json:"quantity"`
	ProductType     string          `json:"productType"`
	SecurityID      string          `json:"securityId"`
	Price           decimal.Decimal `json:"price"`
	TriggerPrice    decimal.Decimal `json:"triggerPrice,omitempty"`
}


type KnowYourMarginResponse struct {
	TotalMargin         decimal.Decimal `json:"totalMargin"`
	SpanMargin          decimal.Decimal `json:"spanMargin"`
	ExposureMargin      decimal.Decimal `json:"exposureMargin"`
	AvailableBalance    decimal.Decimal `json:"availableBalance"`
	VariableMargin      decimal.Decimal `json:"variableMargin"`
	InsufficientBalance decimal.Decimal `json:"insufficientBalance"`
	Brokerage           decimal.Decimal `json:"brokerage"`
	Leverage            string          `json:"leverage"`
}


type MultiScripMarginCalcRequest struct {
	DhanClientID string              `json:"dhanClientId,omitempty"`
	Orders       []KnowYourMarginReq `json:"orders"`
}


type MultiScripMarginCalcResponse struct {
	TotalMargin    decimal.Decimal `json:"totalMargin"`
	SpanMargin     decimal.Decimal `json:"spanMargin"`
	ExposureMargin decimal.Decimal `json:"exposureMargin"`
}


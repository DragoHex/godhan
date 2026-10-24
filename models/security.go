package models

import (
	"github.com/shopspring/decimal"
)

type UserIPRequest struct {
	DhanClientID string `json:"dhanClientId,omitempty"`
	IP           string `json:"ip"`
	IPFlag       string `json:"ipFlag"` // PRIMARY or SECONDARY
}


type UserIPResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}


type GetIPDetailsResponse struct {
	ModifyDatePrimary   string `json:"modifyDatePrimary"`
	ModifyDateSecondary string `json:"modifyDateSecondary"`
	PrimaryIP           string `json:"primaryIP"`
	SecondaryIP         string `json:"secondaryIP"`
	DetectedIP          string `json:"detectedIP"`
	IPMatchStatus       string `json:"ipMatchStatus"`
	OrdersAllowed       bool   `json:"ordersAllowed"`
}


type KillSwitchResponse struct {
	DhanClientID     string `json:"dhanClientId"`
	KillSwitchStatus string `json:"killSwitchStatus"` // ACTIVATE or DEACTIVATE
}


type PnlBasedExitRequest struct {
	DhanClientID     string          `json:"dhanClientId,omitempty"`
	ProfitValue      decimal.Decimal `json:"profitValue"`
	LossValue        decimal.Decimal `json:"lossValue"`
	EnableKillSwitch bool            `json:"enableKillSwitch"`
	ProductType      []string        `json:"productType"`
}


type PnlExitResponse struct {
	Message       string `json:"message"`
	PnlExitStatus string `json:"pnlExitStatus"`
}


type ExitPnlResponse struct {
	PnlExitStatus    string          `json:"pnlExitStatus"`
	Profit           decimal.Decimal `json:"profit"`
	Loss             decimal.Decimal `json:"loss"`
	EnableKillSwitch bool            `json:"enableKillSwitch"`
	ProductType      []string        `json:"productType"`
}


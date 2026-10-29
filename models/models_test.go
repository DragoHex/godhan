package models_test

import (
	"encoding/json"
	"testing"

	"github.com/DragoHex/godhan/models"
	"github.com/shopspring/decimal"
)

func TestDhanOrderModels_Serialization(t *testing.T) {
	orderJSON := `{
		"dhanClientId": "1000000001",
		"orderId": "112111182198",
		"correlationId": "tag123",
		"orderStatus": "CONFIRMED",
		"transactionType": "BUY",
		"exchangeSegment": "NSE_EQ",
		"productType": "INTRADAY",
		"orderType": "LIMIT",
		"validity": "DAY",
		"tradingSymbol": "TATASTEEL",
		"securityId": "11536",
		"quantity": 100,
		"disclosedQuantity": 0,
		"price": 120.50,
		"triggerPrice": 0.0,
		"afterMarketOrder": false,
		"boProfitValue": 0.0,
		"boStopLossValue": 0.0,
		"legName": "ENTRY_LEG",
		"createTime": "2023-10-15 14:30:45",
		"updateTime": "2023-10-15 14:30:45",
		"exchangeTime": "2023-10-15 14:30:45",
		"drvExpiryDate": "",
		"drvOptionType": "",
		"drvStrikePrice": 0.0,
		"omsErrorCode": "",
		"omsErrorDescription": "",
		"filled_qty": 100,
		"costPrice": 120.50
	}`

	var o models.OrderResponse
	if err := json.Unmarshal([]byte(orderJSON), &o); err != nil {
		t.Fatalf("Unmarshal OrderResponse error: %v", err)
	}

	if o.OrderID != "112111182198" || o.TradingSymbol != "TATASTEEL" {
		t.Errorf("Unexpected order fields: %+v", o)
	}
	if !o.Price.Equal(decimal.NewFromFloat(120.50)) {
		t.Errorf("Expected price 120.50, got %s", o.Price)
	}
}

func TestDhanPortfolioModels_Serialization(t *testing.T) {
	posJSON := `{
		"dhanClientId": "1000000001",
		"tradingSymbol": "TATASTEEL",
		"securityId": "11536",
		"positionType": "LONG",
		"exchangeSegment": "NSE_EQ",
		"productType": "INTRADAY",
		"buyQty": 100,
		"costPrice": 120.50,
		"buyAvg": 120.50,
		"sellQty": 0,
		"sellAvg": 0.0,
		"netQty": 100,
		"realizedProfit": 0.0,
		"unrealizedProfit": 50.0,
		"rbiReferenceRate": 1.0,
		"multiplier": 1,
		"carryForwardBuyQty": 0,
		"carryForwardSellQty": 0,
		"carryForwardBuyValue": 0.0,
		"carryForwardSellValue": 0.0,
		"dayBuyQty": 100,
		"daySellQty": 0,
		"dayBuyValue": 12050.0,
		"daySellValue": 0.0,
		"drvExpiryDate": "",
		"drvOptionType": "",
		"drvStrikePrice": 0.0,
		"crossCurrency": false
	}`

	var p models.Position
	if err := json.Unmarshal([]byte(posJSON), &p); err != nil {
		t.Fatalf("Unmarshal Position error: %v", err)
	}

	if p.TradingSymbol != "TATASTEEL" || p.NetQty != 100 {
		t.Errorf("Unexpected position: %+v", p)
	}
}

func TestDhanPostbackPayload_Serialization(t *testing.T) {
	pbJSON := `{
		"dhanClientId": "1000000001",
		"orderId": "112111182198",
		"correlationId": "tag123",
		"orderStatus": "TRADED",
		"transactionType": "BUY",
		"exchangeSegment": "NSE_EQ",
		"productType": "INTRADAY",
		"orderType": "LIMIT",
		"validity": "DAY",
		"tradingSymbol": "TATASTEEL",
		"securityId": "11536",
		"quantity": 100,
		"price": 120.50,
		"triggerPrice": 0.0,
		"createTime": "2023-10-15 14:30:45",
		"updateTime": "2023-10-15 14:30:45",
		"exchangeTime": "2023-10-15 14:30:45",
		"exchangeOrderId": "2023101500001"
	}`

	var pb models.PostbackPayload
	if err := json.Unmarshal([]byte(pbJSON), &pb); err != nil {
		t.Fatalf("Unmarshal PostbackPayload error: %v", err)
	}

	if pb.OrderID != "112111182198" || pb.OrderStatus != "TRADED" {
		t.Errorf("Unexpected postback payload: %+v", pb)
	}
}

func TestRenewTokenResponse_Serialization(t *testing.T) {
	rawJSON := `{
		"dhanClientId": "1000000001",
		"dhanClientName": "Jane Doe",
		"dhanClientUcc": "UCC999",
		"givenPowerOfAttorney": true,
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		"expiryTime": "2026-10-27 15:30:00"
	}`

	var resp models.RenewTokenResponse
	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		t.Fatalf("Unmarshal RenewTokenResponse failed: %v", err)
	}

	if resp.DhanClientID != "1000000001" {
		t.Errorf("DhanClientID = %q, want '1000000001'", resp.DhanClientID)
	}
	if resp.DhanClientName != "Jane Doe" {
		t.Errorf("DhanClientName = %q, want 'Jane Doe'", resp.DhanClientName)
	}
	if resp.DhanClientUCC != "UCC999" {
		t.Errorf("DhanClientUCC = %q, want 'UCC999'", resp.DhanClientUCC)
	}
	if !resp.GivenPowerOfAttorney {
		t.Errorf("GivenPowerOfAttorney = false, want true")
	}
	if resp.AccessToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test" {
		t.Errorf("AccessToken = %q, want expected JWT", resp.AccessToken)
	}
	if resp.ExpiryTime.Time.Year() != 2026 {
		t.Errorf("ExpiryTime year = %d, want 2026", resp.ExpiryTime.Time.Year())
	}

	marshaled, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal RenewTokenResponse failed: %v", err)
	}

	var roundTrip models.RenewTokenResponse
	if err := json.Unmarshal(marshaled, &roundTrip); err != nil {
		t.Fatalf("Round-trip Unmarshal failed: %v", err)
	}

	if roundTrip.AccessToken != resp.AccessToken {
		t.Errorf("roundTrip.AccessToken = %q, want %q", roundTrip.AccessToken, resp.AccessToken)
	}
}

func TestRenewTokenResponse_LiveGatewayFormat(t *testing.T) {
	rawJSON := `{
		"createTime": "2026-10-29T15:08:58.164",
		"expiryTime": "2026-10-30T15:08:58.163",
		"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.live-token-value"
	}`

	var resp models.RenewTokenResponse
	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		t.Fatalf("Unmarshal live gateway RenewTokenResponse failed: %v", err)
	}

	if resp.Token != "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.live-token-value" {
		t.Errorf("Token = %q, want expected token", resp.Token)
	}
	if resp.AccessToken != "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.live-token-value" {
		t.Errorf("AccessToken = %q, want expected token mapped to AccessToken", resp.AccessToken)
	}
	if resp.CreateTime.Time.IsZero() {
		t.Errorf("expected non-zero CreateTime")
	}
	if resp.ExpiryTime.Time.IsZero() {
		t.Errorf("expected non-zero ExpiryTime")
	}
}



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

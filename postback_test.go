package godhan

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestParsePostback_Success(t *testing.T) {
	rawJSON := `{
		"dhanClientId": "1000000001",
		"orderId": "112111182198",
		"correlationId": "CLIENT_TAG_001",
		"orderStatus": "TRADED",
		"transactionType": "BUY",
		"exchangeSegment": "NSE_EQ",
		"productType": "CNC",
		"orderType": "LIMIT",
		"validity": "DAY",
		"tradingSymbol": "TCS",
		"securityId": "11536",
		"quantity": 10,
		"disclosedQuantity": 0,
		"price": 3500.50,
		"triggerPrice": 0.00,
		"afterMarketOrder": false,
		"remainingQuantity": 0,
		"filledQty": 10,
		"averageTradedPrice": 3500.50,
		"createTime": "2025-01-16 09:15:00",
		"updateTime": "2025-01-16 09:15:02",
		"exchangeTime": "2025-01-16 09:15:02",
		"omsErrorCode": "0",
		"omsErrorDescription": "SUCCESS"
	}`

	payload, err := ParsePostback([]byte(rawJSON))
	if err != nil {
		t.Fatalf("unexpected error parsing postback: %v", err)
	}

	if payload.DhanClientID != "1000000001" {
		t.Errorf("DhanClientID = %q, want 1000000001", payload.DhanClientID)
	}
	if payload.OrderID != "112111182198" {
		t.Errorf("OrderID = %q, want 112111182198", payload.OrderID)
	}
	if payload.CorrelationID != "CLIENT_TAG_001" {
		t.Errorf("CorrelationID = %q, want CLIENT_TAG_001", payload.CorrelationID)
	}
	if payload.OrderStatus != "TRADED" {
		t.Errorf("OrderStatus = %q, want TRADED", payload.OrderStatus)
	}
	if payload.Quantity != 10 || payload.FilledQty != 10 {
		t.Errorf("Quantity = %d, FilledQty = %d, want 10", payload.Quantity, payload.FilledQty)
	}
	if !payload.Price.Equal(decimal.NewFromFloat(3500.50)) {
		t.Errorf("Price = %v, want 3500.50", payload.Price)
	}
	if !payload.AverageTradedPrice.Equal(decimal.NewFromFloat(3500.50)) {
		t.Errorf("AverageTradedPrice = %v, want 3500.50", payload.AverageTradedPrice)
	}
}

func TestParsePostbackRequest(t *testing.T) {
	rawJSON := `{"orderId": "112111182198", "orderStatus": "REJECTED", "dhanClientId": "1000000001"}`
	req := httptest.NewRequest(http.MethodPost, "/broker-callback", bytes.NewReader([]byte(rawJSON)))

	payload, err := ParsePostbackRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload.OrderID != "112111182198" || payload.OrderStatus != "REJECTED" {
		t.Errorf("unexpected payload: %+v", payload)
	}

	// Malformed body
	badReq := httptest.NewRequest(http.MethodPost, "/broker-callback", bytes.NewReader([]byte("not-json")))
	_, err = ParsePostbackRequest(badReq)
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for malformed postback request, got %v", err)
	}
}

func TestIsTerminalStatus(t *testing.T) {
	terminals := []string{"TRADED", "REJECTED", "CANCELLED", "EXPIRED"}
	for _, s := range terminals {
		if !IsTerminalStatus(s) {
			t.Errorf("IsTerminalStatus(%q) = false, want true", s)
		}
	}

	nonTerminals := []string{"TRANSIT", "PENDING", "PARTIALLY_FILLED"}
	for _, s := range nonTerminals {
		if IsTerminalStatus(s) {
			t.Errorf("IsTerminalStatus(%q) = true, want false", s)
		}
	}
}

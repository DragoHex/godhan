package godhan

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrders_ExitOrder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/orders/112111182198" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"orderId": "112111182198",
			"orderStatus": "CANCELLED"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "valid-token", WithBaseURI(ts.URL))

	resp, err := client.ExitOrder(context.Background(), "112111182198")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.OrderStatus != "CANCELLED" {
		t.Errorf("OrderStatus = %q, want CANCELLED", resp.OrderStatus)
	}
}

func TestOrders_GetOrderHistory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders/112111182198" {
			t.Errorf("path = %q, want /orders/112111182198", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"orderId": "112111182198",
			"orderStatus": "TRADED",
			"tradingSymbol": "TCS"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "valid-token", WithBaseURI(ts.URL))

	resp, err := client.GetOrderHistory(context.Background(), "112111182198")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.OrderID != "112111182198" || resp.OrderStatus != "TRADED" {
		t.Errorf("unexpected resp: %+v", resp)
	}
}

func TestTrades_GetOrderTrades(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/trades/112111182198" {
			t.Errorf("path = %q, want /trades/112111182198", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{
				"orderId": "112111182198",
				"exchangeTradeId": "TRD_111",
				"tradingSymbol": "TCS",
				"tradedQuantity": 10
			}
		]`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "valid-token", WithBaseURI(ts.URL))

	trades, err := client.GetOrderTrades(context.Background(), "112111182198")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(trades) != 1 || trades[0].ExchangeTradeID != "TRD_111" {
		t.Errorf("unexpected trades: %+v", trades)
	}
}

package godhan

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrades_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == URITrades:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"orderId": "112233",
					"exchangeTradeId": "TR-9900",
					"tradingSymbol": "RELIANCE",
					"tradedQuantity": 50,
					"tradedPrice": 2450.00
				}
			]`))
		case r.Method == http.MethodGet && r.URL.Path == URITrades+"/112233":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"orderId": "112233",
					"exchangeTradeId": "TR-9900",
					"tradedQuantity": 50,
					"tradedPrice": 2450.00
				}
			]`))
		case r.Method == http.MethodGet && r.URL.Path == URITrades+"/2023-10-01/2023-10-15/1":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"orderId": "hist-1",
					"exchangeTradeId": "TR-HIST-1",
					"tradedQuantity": 100,
					"tradedPrice": 1200.00
				}
			]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GetTrades
	trades, err := client.GetTrades(ctx)
	if err != nil {
		t.Fatalf("GetTrades failed: %v", err)
	}
	if len(trades) != 1 || trades[0].ExchangeTradeID != "TR-9900" {
		t.Errorf("unexpected trades response: %+v", trades)
	}

	// 2. GetTradesByOrderID
	orderTrades, err := client.GetTradesByOrderID(ctx, "112233")
	if err != nil {
		t.Fatalf("GetTradesByOrderID failed: %v", err)
	}
	if len(orderTrades) != 1 || orderTrades[0].OrderID != "112233" {
		t.Errorf("unexpected order trades: %+v", orderTrades)
	}

	// 3. GetTradeHistory
	histTrades, err := client.GetTradeHistory(ctx, "2023-10-01", "2023-10-15", 1)
	if err != nil {
		t.Fatalf("GetTradeHistory failed: %v", err)
	}
	if len(histTrades) != 1 || histTrades[0].ExchangeTradeID != "TR-HIST-1" {
		t.Errorf("unexpected hist trades: %+v", histTrades)
	}
}

func TestTrades_Validation(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Empty models.Order ID
	_, err := client.GetTradesByOrderID(ctx, "")
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}

	// Invalid page
	_, err = client.GetTradeHistory(ctx, "2023-10-01", "2023-10-15", 0)
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for page 0, got %v", err)
	}
}

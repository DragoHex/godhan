package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestOrders_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == URIOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"1122334455","orderStatus":"TRANSIT"}`))
		case r.Method == http.MethodPut && r.URL.Path == URIOrders+"/1122334455":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"1122334455","orderStatus":"PENDING"}`))
		case r.Method == http.MethodDelete && r.URL.Path == URIOrders+"/1122334455":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"1122334455","orderStatus":"CANCELLED"}`))
		case r.Method == http.MethodGet && r.URL.Path == URIOrders+"/1122334455":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"orderId": "1122334455",
				"dhanClientId": "test-client",
				"orderStatus": "PENDING",
				"transactionType": "BUY",
				"exchangeSegment": "NSE_EQ",
				"productType": "CNC",
				"orderType": "LIMIT",
				"tradingSymbol": "INFY",
				"securityId": "1594",
				"quantity": 100,
				"price": 1450.50
			}`))
		case r.Method == http.MethodGet && r.URL.Path == URIOrdersExternal+"/corr-999":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"orderId": "1122334455",
				"correlationId": "corr-999",
				"orderStatus": "PENDING"
			}`))
		case r.Method == http.MethodGet && r.URL.Path == URIOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"orderId":"1122334455","orderStatus":"PENDING"},
				{"orderId":"9988776655","orderStatus":"COMPLETE"}
			]`))
		case r.Method == http.MethodPost && r.URL.Path == URIOrdersSlicing:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"orderId":"slice-1","orderStatus":"TRANSIT"},
				{"orderId":"slice-2","orderStatus":"TRANSIT"}
			]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. PlaceOrder
	placeResp, err := client.PlaceOrder(ctx, models.OrderRequest{
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		ProductType:     "CNC",
		OrderType:       "LIMIT",
		Validity:        "DAY",
		TradingSymbol:   "INFY",
		SecurityID:      "1594",
		Quantity:        100,
		Price:           decimal.NewFromFloat(1450.50),
	})
	if err != nil {
		t.Fatalf("PlaceOrder failed: %v", err)
	}
	if placeResp.OrderID != "1122334455" || placeResp.OrderStatus != "TRANSIT" {
		t.Errorf("unexpected place order response: %+v", placeResp)
	}

	// 2. ModifyOrder
	modResp, err := client.ModifyOrder(ctx, "1122334455", models.OrderModifyRequest{
		OrderType: "LIMIT",
		Price:     decimal.NewFromFloat(1455.00),
		Quantity:  100,
		Validity:  "DAY",
	})
	if err != nil {
		t.Fatalf("ModifyOrder failed: %v", err)
	}
	if modResp.OrderID != "1122334455" || modResp.OrderStatus != "PENDING" {
		t.Errorf("unexpected modify order response: %+v", modResp)
	}

	// 3. CancelOrder
	cancelResp, err := client.CancelOrder(ctx, "1122334455")
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}
	if cancelResp.OrderStatus != "CANCELLED" {
		t.Errorf("expected CANCELLED, got %s", cancelResp.OrderStatus)
	}

	// 4. GetOrderByID
	order, err := client.GetOrderByID(ctx, "1122334455")
	if err != nil {
		t.Fatalf("GetOrderByID failed: %v", err)
	}
	if order.OrderID != "1122334455" || order.Quantity != 100 {
		t.Errorf("unexpected order details: %+v", order)
	}

	// 5. GetOrderByCorrelationID
	corrOrder, err := client.GetOrderByCorrelationID(ctx, "corr-999")
	if err != nil {
		t.Fatalf("GetOrderByCorrelationID failed: %v", err)
	}
	if corrOrder.CorrelationID != "corr-999" {
		t.Errorf("unexpected correlation order: %+v", corrOrder)
	}

	// 6. GetOrders (models.Order Book)
	orders, err := client.GetOrders(ctx)
	if err != nil {
		t.Fatalf("GetOrders failed: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}

	// 7. SliceOrder
	sliceResp, err := client.SliceOrder(ctx, models.OrderRequest{
		TransactionType: "BUY",
		ExchangeSegment: "NSE_FNO",
		ProductType:     "INTRADAY",
		OrderType:       "MARKET",
		Validity:        "DAY",
		SecurityID:      "52175",
		Quantity:        5000,
	})
	if err != nil {
		t.Fatalf("SliceOrder failed: %v", err)
	}
	if len(sliceResp) != 2 {
		t.Errorf("expected 2 sliced orders, got %d", len(sliceResp))
	}
}

func TestOrders_ValidationNegative(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Zero quantity
	_, err := client.PlaceOrder(ctx, models.OrderRequest{
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		Quantity:        0,
	})
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for zero quantity, got %v", err)
	}

	// Empty OrderID cancel
	_, err = client.CancelOrder(ctx, "")
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for empty order ID, got %v", err)
	}
}

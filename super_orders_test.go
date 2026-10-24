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

func TestSuperOrders_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == URISuperOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"super-101","orderStatus":"TRANSIT"}`))
		case r.Method == http.MethodPut && r.URL.Path == URISuperOrders+"/super-101":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"super-101","orderStatus":"PENDING"}`))
		case r.Method == http.MethodDelete && r.URL.Path == URISuperOrders+"/super-101/STOP_LOSS_LEG":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"super-101","orderStatus":"CANCELLED"}`))
		case r.Method == http.MethodGet && r.URL.Path == URISuperOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"orderId": "super-101",
					"dhanClientId": "test-client",
					"orderStatus": "PENDING",
					"transactionType": "BUY",
					"price": 100.0,
					"targetPrice": 110.0,
					"stopLossPrice": 95.0
				}
			]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. PlaceSuperOrder
	placeResp, err := client.PlaceSuperOrder(ctx, models.SuperOrderRequest{
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		ProductType:     "INTRADAY",
		OrderType:       "LIMIT",
		SecurityID:      "1333",
		Quantity:        10,
		Price:           decimal.NewFromFloat(100.0),
		TargetPrice:     decimal.NewFromFloat(110.0),
		StopLossPrice:   decimal.NewFromFloat(95.0),
		TrailingJump:    decimal.NewFromFloat(1.0),
	})
	if err != nil {
		t.Fatalf("PlaceSuperOrder failed: %v", err)
	}
	if placeResp.OrderID != "super-101" {
		t.Errorf("unexpected super order response: %+v", placeResp)
	}

	// 2. ModifySuperOrder
	modResp, err := client.ModifySuperOrder(ctx, "super-101", models.SuperModifyRequest{
		OrderType:     "LIMIT",
		LegName:       "TARGET_LEG",
		Quantity:      10,
		Price:         decimal.NewFromFloat(100.0),
		TargetPrice:   decimal.NewFromFloat(115.0),
		StopLossPrice: decimal.NewFromFloat(95.0),
	})
	if err != nil {
		t.Fatalf("ModifySuperOrder failed: %v", err)
	}
	if modResp.OrderID != "super-101" {
		t.Errorf("unexpected modify super order response: %+v", modResp)
	}

	// 3. CancelSuperOrder
	cancelResp, err := client.CancelSuperOrder(ctx, "super-101", "STOP_LOSS_LEG")
	if err != nil {
		t.Fatalf("CancelSuperOrder failed: %v", err)
	}
	if cancelResp.OrderStatus != "CANCELLED" {
		t.Errorf("expected CANCELLED, got %s", cancelResp.OrderStatus)
	}

	// 4. GetSuperOrders
	superOrders, err := client.GetSuperOrders(ctx)
	if err != nil {
		t.Fatalf("GetSuperOrders failed: %v", err)
	}
	if len(superOrders) != 1 || superOrders[0].OrderID != "super-101" {
		t.Errorf("unexpected super orders: %+v", superOrders)
	}
}

func TestSuperOrders_ValidationNegative(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Zero quantity
	_, err := client.PlaceSuperOrder(ctx, models.SuperOrderRequest{
		TransactionType: "BUY",
		Quantity:        0,
	})
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}

	// Empty leg on cancel
	_, err = client.CancelSuperOrder(ctx, "super-101", "")
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for empty leg, got %v", err)
	}
}

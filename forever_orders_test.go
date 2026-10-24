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

func TestForeverOrders_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == URIForeverOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"gtt-555","orderStatus":"PENDING"}`))
		case r.Method == http.MethodPut && r.URL.Path == URIForeverOrders+"/gtt-555":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"gtt-555","orderStatus":"PENDING"}`))
		case r.Method == http.MethodDelete && r.URL.Path == URIForeverOrders+"/gtt-555":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"orderId":"gtt-555","orderStatus":"CANCELLED"}`))
		case r.Method == http.MethodGet && r.URL.Path == URIForeverOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"orderId":"gtt-555","orderStatus":"PENDING"}]`))
		case r.Method == http.MethodPost && r.URL.Path == URIAlertOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"alertId":"alert-101","alertStatus":"ACTIVE"}`))
		case r.Method == http.MethodGet && r.URL.Path == URIAlertOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"alertId":"alert-101","alertStatus":"ACTIVE"}]`))
		case r.Method == http.MethodGet && r.URL.Path == URIAlertOrders+"/alert-101":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"alertId":"alert-101","alertStatus":"ACTIVE"}`))
		case r.Method == http.MethodPut && r.URL.Path == URIAlertOrders+"/alert-101":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"alertId":"alert-101","alertStatus":"ACTIVE"}`))
		case r.Method == http.MethodDelete && r.URL.Path == URIAlertOrders+"/alert-101":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"alertId":"alert-101","alertStatus":"CANCELLED"}`))
		case r.Method == http.MethodPost && r.URL.Path == URIMultiAlertOrders:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"multiOrderId":"multi-77","status":"SUCCESS"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. PlaceForeverOrder
	placeResp, err := client.PlaceForeverOrder(ctx, models.GTTOrderModel{
		OrderFlag:       "SINGLE",
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		ProductType:     "CNC",
		OrderType:       "LIMIT",
		SecurityID:      "1333",
		Quantity:        50,
		Price:           decimal.NewFromFloat(150.0),
		TriggerPrice:    decimal.NewFromFloat(152.0),
	})
	if err != nil {
		t.Fatalf("PlaceForeverOrder failed: %v", err)
	}
	if placeResp.OrderID != "gtt-555" {
		t.Errorf("unexpected gtt response: %+v", placeResp)
	}

	// 2. ModifyForeverOrder
	modResp, err := client.ModifyForeverOrder(ctx, "gtt-555", models.GttModifyRequest{
		OrderFlag:    "SINGLE",
		OrderType:    "LIMIT",
		Quantity:     60,
		Price:        decimal.NewFromFloat(155.0),
		TriggerPrice: decimal.NewFromFloat(157.0),
	})
	if err != nil {
		t.Fatalf("ModifyForeverOrder failed: %v", err)
	}
	if modResp.OrderID != "gtt-555" {
		t.Errorf("unexpected modify gtt response: %+v", modResp)
	}

	// 3. CancelForeverOrder
	cancelResp, err := client.CancelForeverOrder(ctx, "gtt-555")
	if err != nil {
		t.Fatalf("CancelForeverOrder failed: %v", err)
	}
	if cancelResp.OrderStatus != "CANCELLED" {
		t.Errorf("expected CANCELLED, got %s", cancelResp.OrderStatus)
	}

	// 4. GetForeverOrders
	gtts, err := client.GetForeverOrders(ctx)
	if err != nil {
		t.Fatalf("GetForeverOrders failed: %v", err)
	}
	if len(gtts) != 1 || gtts[0].OrderID != "gtt-555" {
		t.Errorf("unexpected gtts list: %+v", gtts)
	}

	// 5. Alerts lifecycle
	alertResp, err := client.PlaceAlertOrder(ctx, models.AlertOrderRequest{})
	if err != nil {
		t.Fatalf("PlaceAlertOrder failed: %v", err)
	}
	if alertResp.AlertID != "alert-101" {
		t.Errorf("unexpected alert response: %+v", alertResp)
	}

	alerts, err := client.GetAlertOrders(ctx)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("GetAlertOrders failed: %v", err)
	}

	getAlert, err := client.GetAlertOrderByID(ctx, "alert-101")
	if err != nil || getAlert.AlertID != "alert-101" {
		t.Fatalf("GetAlertOrderByID failed: %v", err)
	}

	modAlert, err := client.ModifyAlertOrder(ctx, "alert-101", models.AlertModifyRequest{})
	if err != nil || modAlert.AlertStatus != "ACTIVE" {
		t.Fatalf("ModifyAlertOrder failed: %v", err)
	}

	delAlert, err := client.CancelAlertOrder(ctx, "alert-101")
	if err != nil || delAlert.AlertStatus != "CANCELLED" {
		t.Fatalf("CancelAlertOrder failed: %v", err)
	}

	// Multi order
	multiResp, err := client.PlaceMultiOrder(ctx, models.MultiOrderRequest{})
	if err != nil || multiResp.MultiOrderID != "multi-77" {
		t.Fatalf("PlaceMultiOrder failed: %v", err)
	}
}

func TestForeverOrders_Validation(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Empty ID on cancel
	_, err := client.CancelForeverOrder(ctx, "")
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

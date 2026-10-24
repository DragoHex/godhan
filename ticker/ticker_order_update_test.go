package ticker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"go.uber.org/goleak"
)

var testUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func TestTicker_OnOrderUpdate_Callback(t *testing.T) {
	defer goleak.VerifyNone(t)

	var orderReceived sync.WaitGroup
	orderReceived.Add(1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// Send order execution event frame
		orderFrame := `{
			"type": "order",
			"data": {
				"dhanClientId": "1000000001",
				"orderId": "112111182198",
				"correlationId": "CLIENT_TAG_001",
				"orderStatus": "TRADED",
				"transactionType": "BUY",
				"exchangeSegment": "NSE_EQ",
				"productType": "CNC",
				"tradingSymbol": "TCS",
				"securityId": "11536",
				"quantity": 10,
				"filledQty": 10,
				"averageTradedPrice": 3500.50
			}
		}`
		_ = conn.WriteMessage(websocket.TextMessage, []byte(orderFrame))

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	ticker := NewTicker("1000000001", "test-token",
		WithTickerURI(wsURL),
		WithAutoReconnect(false),
	)

	var receivedOrder OrderUpdate
	ticker.OnOrderUpdate(func(ord OrderUpdate) {
		receivedOrder = ord
		orderReceived.Done()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := ticker.Connect(ctx); err != nil {
		t.Fatalf("failed to connect ticker: %v", err)
	}

	orderReceived.Wait()

	if receivedOrder.OrderID != "112111182198" {
		t.Errorf("OrderID = %q, want 112111182198", receivedOrder.OrderID)
	}
	if receivedOrder.OrderStatus != "TRADED" {
		t.Errorf("OrderStatus = %q, want TRADED", receivedOrder.OrderStatus)
	}
	if !receivedOrder.AverageTradedPrice.Equal(decimal.NewFromFloat(3500.50)) {
		t.Errorf("AverageTradedPrice = %v, want 3500.50", receivedOrder.AverageTradedPrice)
	}

	if err := ticker.Close(); err != nil {
		t.Fatalf("failed to close ticker: %v", err)
	}
}

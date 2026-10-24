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
	"go.uber.org/goleak"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func TestTicker_LifecycleAndStreaming(t *testing.T) {
	defer goleak.VerifyNone(t)

	var subReceived sync.WaitGroup
	subReceived.Add(1)

	// Mock WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Send initial authorization/welcome or handle subscription message
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			// When client sends subscribe request
			if len(msg) > 0 {
				subReceived.Done()
				// Send a binary LTP tick back
				tickBytes := buildBinaryPacket(11, 1, 11536, []byte{
					0x00, 0x00, 0xb5, 0x44, // float32 1448.0
					0x00, 0x10, 0x2b, 0x65, // timestamp
				})
				_ = conn.WriteMessage(websocket.BinaryMessage, tickBytes)
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	ticker := NewTicker("test-client", "test-token",
		WithTickerURI(wsURL),
		WithAutoReconnect(false),
	)

	var receivedTick Tick
	tickChan := make(chan Tick, 1)

	ticker.OnTick(func(tick Tick) {
		tickChan <- tick
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := ticker.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !ticker.IsConnected() {
		t.Fatal("expected ticker.IsConnected() == true")
	}

	// Subscribe
	err := ticker.Subscribe([]TickerSubscription{
		{ExchangeSegment: 1, SecurityID: 11536, Mode: TickerModeLTP},
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Wait for tick or timeout
	select {
	case receivedTick = <-tickChan:
		if receivedTick.SecurityID != 11536 {
			t.Errorf("expected security ID 11536, got %d", receivedTick.SecurityID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for binary tick")
	}

	// Close cleanly
	if err := ticker.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if ticker.IsConnected() {
		t.Fatal("expected ticker.IsConnected() == false after Close")
	}
}

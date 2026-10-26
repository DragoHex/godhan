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

func TestTicker_WithTickerProxyURL(t *testing.T) {
	// 1. Option setting verification
	ticker := NewTicker("test-client", "test-token",
		WithTickerProxyURL("http://proxy.local:8080"),
	)
	if ticker.proxyURL == nil || ticker.proxyURL.String() != "http://proxy.local:8080" {
		t.Fatalf("expected proxyURL http://proxy.local:8080, got %v", ticker.proxyURL)
	}

	// 2. Whitespace / empty proxy URL handling (should not set nil proxyURL)
	tickerEmpty := NewTicker("test-client", "test-token",
		WithTickerProxyURL("   "),
	)
	if tickerEmpty.proxyURL != nil {
		t.Errorf("expected nil proxyURL for empty string, got %v", tickerEmpty.proxyURL)
	}

	// 3. Connect through unreachable proxy fails with ErrBrokerUnavailable
	unreachableTicker := NewTicker("test-client", "test-token",
		WithTickerURI("ws://127.0.0.1:59999/feed"),
		WithTickerProxyURL("http://127.0.0.1:59998"),
		WithAutoReconnect(false),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := unreachableTicker.Connect(ctx)
	if err == nil {
		t.Fatal("expected connection error through unreachable proxy, got nil")
	}
}

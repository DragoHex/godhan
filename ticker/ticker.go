package ticker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

// OrderUpdate represents an order update event pushed via WebSocket

const DefaultFeedURI = "wss://feed.dhan.co"

var (
	ErrBrokerUnavailable = errors.New("broker streaming gateway unreachable")
	ErrInvalidInput      = errors.New("invalid input or packet")
)

type OrderUpdate struct {
	DhanClientID        string          `json:"dhanClientId"`
	OrderID             string          `json:"orderId"`
	CorrelationID       string          `json:"correlationId"`
	OrderStatus         string          `json:"orderStatus"`
	TransactionType     string          `json:"transactionType"`
	ExchangeSegment     string          `json:"exchangeSegment"`
	ProductType         string          `json:"productType"`
	OrderType           string          `json:"orderType,omitempty"`
	Validity            string          `json:"validity,omitempty"`
	TradingSymbol       string          `json:"tradingSymbol"`
	SecurityID          string          `json:"securityId"`
	Quantity            int32           `json:"quantity"`
	DisclosedQuantity   int32           `json:"disclosedQuantity,omitempty"`
	Price               decimal.Decimal `json:"price,omitempty"`
	TriggerPrice        decimal.Decimal `json:"triggerPrice,omitempty"`
	AfterMarketOrder    bool            `json:"afterMarketOrder,omitempty"`
	RemainingQuantity   int32           `json:"remainingQuantity,omitempty"`
	FilledQty           int32           `json:"filledQty"`
	AverageTradedPrice  decimal.Decimal `json:"averageTradedPrice"`
	CreateTime          string          `json:"createTime,omitempty"`
	UpdateTime          string          `json:"updateTime,omitempty"`
	ExchangeTime        string          `json:"exchangeTime,omitempty"`
	OmsErrorCode        string          `json:"omsErrorCode,omitempty"`
	OmsErrorDescription string          `json:"omsErrorDescription,omitempty"`
}

type wsOrderEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TickerSubscription defines an instrument to stream
type TickerSubscription struct {
	ExchangeSegment uint8  `json:"exchangeSegment"`
	SecurityID      uint32 `json:"securityId"`
	Mode            int    `json:"mode"` // TickerModeLTP, TickerModeQuote, TickerModeDepth
}

// TickerOption configures Ticker behavior
type TickerOption func(*Ticker)

// WithTickerURI overrides WebSocket endpoint
func WithTickerURI(uri string) TickerOption {
	return func(t *Ticker) {
		if uri != "" {
			t.uri = uri
		}
	}
}

// WithAutoReconnect controls automatic reconnection on unexpected disconnection
func WithAutoReconnect(enabled bool) TickerOption {
	return func(t *Ticker) {
		t.autoReconnect = enabled
	}
}

// WithTickerProxyURL configures a forward proxy for WebSocket dialer.
func WithTickerProxyURL(rawURL string) TickerOption {
	return func(t *Ticker) {
		trimmed := strings.TrimSpace(rawURL)
		if trimmed != "" {
			if u, err := url.Parse(trimmed); err == nil && u.Host != "" {
				t.proxyURL = u
			}
		}
	}
}

// Ticker represents persistent WebSocket streaming client for live market quotes
type Ticker struct {
	clientID      string
	accessToken   string
	uri           string
	autoReconnect bool
	proxyURL      *url.URL

	conn        *websocket.Conn
	connMu      sync.Mutex
	isConnected atomic.Bool
	closed      atomic.Bool

	subsMu sync.RWMutex
	subs   map[string]TickerSubscription

	lastPing atomic.Value // time.Time

	// Callbacks
	onTick        func(Tick)
	onOrderUpdate func(OrderUpdate)
	onError       func(error)
	onConnect     func()
	onClose       func(int, string)
	onReconnect   func(int, time.Duration)

	// Goroutine management
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewTicker constructs a new WebSocket ticker client
func NewTicker(clientID, accessToken string, opts ...TickerOption) *Ticker {
	t := &Ticker{
		clientID:      clientID,
		accessToken:   accessToken,
		uri:           DefaultFeedURI,
		autoReconnect: true,
		subs:          make(map[string]TickerSubscription),
		stopCh:        make(chan struct{}),
	}
	t.lastPing.Store(time.Now())

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// OnTick registers callback for incoming market ticks
func (t *Ticker) OnTick(cb func(Tick)) { t.onTick = cb }

// OnOrderUpdate registers callback for real-time order update events
func (t *Ticker) OnOrderUpdate(cb func(OrderUpdate)) { t.onOrderUpdate = cb }

// OnError registers callback for streaming errors
func (t *Ticker) OnError(cb func(error)) { t.onError = cb }

// OnConnect registers callback triggered upon connection
func (t *Ticker) OnConnect(cb func()) { t.onConnect = cb }

// OnClose registers callback triggered upon disconnection
func (t *Ticker) OnClose(cb func(int, string)) { t.onClose = cb }

// OnReconnect registers callback triggered when auto-reconnecting
func (t *Ticker) OnReconnect(cb func(int, time.Duration)) { t.onReconnect = cb }

// IsConnected returns current connection state
func (t *Ticker) IsConnected() bool {
	return t.isConnected.Load()
}

// Connect dials the WebSocket feed gateway and starts read/watchdog loops
func (t *Ticker) Connect(ctx context.Context) error {
	t.connMu.Lock()
	defer t.connMu.Unlock()

	if t.isConnected.Load() {
		return nil
	}

	headers := http.Header{}
	headers.Set("client-id", t.clientID)
	headers.Set("access-token", t.accessToken)

	dialer := websocket.DefaultDialer
	if t.proxyURL != nil {
		dialer = &websocket.Dialer{
			Proxy:            http.ProxyURL(t.proxyURL),
			HandshakeTimeout: 10 * time.Second,
		}
	}
	conn, _, err := dialer.DialContext(ctx, t.uri, headers)
	if err != nil {
		return fmt.Errorf("%w: failed to dial websocket: %v", ErrBrokerUnavailable, err)
	}

	t.conn = conn
	t.isConnected.Store(true)
	t.closed.Store(false)
	t.lastPing.Store(time.Now())
	t.stopCh = make(chan struct{})

	if t.onConnect != nil {
		go t.onConnect()
	}

	t.wg.Add(1)
	go t.readLoop()

	t.wg.Add(1)
	go t.watchdogLoop()

	return nil
}

// Close gracefully terminates the connection and releases goroutines
func (t *Ticker) Close() error {
	if t.closed.Swap(true) {
		return nil
	}

	close(t.stopCh)

	t.connMu.Lock()
	if t.conn != nil {
		_ = t.conn.Close()
	}
	t.connMu.Unlock()

	t.isConnected.Store(false)
	t.wg.Wait()

	if t.onClose != nil {
		t.onClose(websocket.CloseNormalClosure, "closed by client")
	}

	return nil
}

// Subscribe registers interest in specified instruments
func (t *Ticker) Subscribe(subs []TickerSubscription) error {
	t.subsMu.Lock()
	for _, s := range subs {
		key := fmt.Sprintf("%d-%d", s.ExchangeSegment, s.SecurityID)
		t.subs[key] = s
	}
	t.subsMu.Unlock()

	return t.sendSubscriptionPayload(15, subs)
}

// Unsubscribe removes interest from instruments
func (t *Ticker) Unsubscribe(subs []TickerSubscription) error {
	t.subsMu.Lock()
	for _, s := range subs {
		key := fmt.Sprintf("%d-%d", s.ExchangeSegment, s.SecurityID)
		delete(t.subs, key)
	}
	t.subsMu.Unlock()

	return t.sendSubscriptionPayload(16, subs)
}

func (t *Ticker) sendSubscriptionPayload(requestCode int, subs []TickerSubscription) error {
	t.connMu.Lock()
	defer t.connMu.Unlock()

	if !t.isConnected.Load() || t.conn == nil {
		return ErrBrokerUnavailable
	}

	reqData := map[string]interface{}{
		"RequestCode":    requestCode,
		"InstrumentList": subs,
	}

	payload, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	return t.conn.WriteMessage(websocket.TextMessage, payload)
}

func (t *Ticker) readLoop() {
	defer t.wg.Done()

	for {
		select {
		case <-t.stopCh:
			return
		default:
		}

		t.connMu.Lock()
		conn := t.conn
		t.connMu.Unlock()

		if conn == nil {
			return
		}

		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			if !t.closed.Load() {
				t.isConnected.Store(false)
				if t.onError != nil {
					t.onError(err)
				}
				if t.autoReconnect {
					go t.reconnect()
				}
			}
			return
		}

		t.lastPing.Store(time.Now())

		if msgType == websocket.BinaryMessage {
			tick, err := ParseBinaryTick(msg)
			if err != nil {
				if t.onError != nil {
					t.onError(err)
				}
				continue
			}
			if t.onTick != nil {
				t.onTick(*tick)
			}
		} else if msgType == websocket.TextMessage {
			var env wsOrderEnvelope
			if err := json.Unmarshal(msg, &env); err == nil && env.Type == "order" && len(env.Data) > 0 {
				var ord OrderUpdate
				if err := json.Unmarshal(env.Data, &ord); err == nil {
					if t.onOrderUpdate != nil {
						t.onOrderUpdate(ord)
					}
					continue
				}
			}

			// Fallback: direct order payload
			var ord OrderUpdate
			if err := json.Unmarshal(msg, &ord); err == nil && ord.OrderID != "" {
				if t.onOrderUpdate != nil {
					t.onOrderUpdate(ord)
				}
				continue
			}
		}
	}
}

func (t *Ticker) watchdogLoop() {
	defer t.wg.Done()
	ticker := time.NewTicker(2500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			last := t.lastPing.Load().(time.Time)
			if time.Since(last) > 10*time.Second && t.isConnected.Load() && !t.closed.Load() {
				// Socket stale / dead TCP
				t.connMu.Lock()
				if t.conn != nil {
					_ = t.conn.Close()
				}
				t.connMu.Unlock()
				t.isConnected.Store(false)
				return
			}
		}
	}
}

func (t *Ticker) reconnect() {
	if t.closed.Load() {
		return
	}

	attempt := 1
	for {
		if t.closed.Load() {
			return
		}

		// Exponential backoff 1s to 30s with jitter
		delay := time.Duration(1<<attempt) * time.Second
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}
		delay += time.Duration(rand.Intn(500)) * time.Millisecond

		if t.onReconnect != nil {
			t.onReconnect(attempt, delay)
		}

		time.Sleep(delay)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := t.Connect(ctx)
		cancel()

		if err == nil {
			// Auto resubscribe cached active instruments
			t.subsMu.RLock()
			var list []TickerSubscription
			for _, s := range t.subs {
				list = append(list, s)
			}
			t.subsMu.RUnlock()

			if len(list) > 0 {
				_ = t.Subscribe(list)
			}
			return
		}

		attempt++
	}
}

// ServeWithContext connects and blocks until context cancellation
func (t *Ticker) ServeWithContext(ctx context.Context) error {
	if err := t.Connect(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	return t.Close()
}

# godhan

[![Go Reference](https://pkg.go.dev/badge/github.com/DragoHex/godhan.svg)](https://pkg.go.dev/github.com/DragoHex/godhan)
[![Go Report Card](https://goreportcard.com/badge/github.com/DragoHex/godhan)](https://goreportcard.com/report/github.com/DragoHex/godhan)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

High-performance, production-grade Go SDK for the **DhanHQ API v2** trading infrastructure. Built with 1:1 functional parity in the image of `gokiteconnect`, featuring strict decimal precision, single-socket multiplexed streaming (live ticks + order postbacks), Forever/GTT orders, super orders, portfolio management, and historical chart parsing.

---

## Architecture & Package Layout

`godhan` enforces clean package boundaries with domain isolation:

```
godhan/
├── client.go              # Root Client & functional configuration options
├── transport.go           # HTTPClient seam, connection pooling & header injection
├── endpoints.go           # DhanHQ API v2 REST URI constants
├── errors.go              # Structured APIError and domain sentinels
├── orders.go              # Regular order execution, modification, slicing & cancellation
├── forever_orders.go      # Forever (GTT / OCO) order lifecycle
├── super_orders.go        # Target & Stop-loss multi-leg Bracket orders
├── portfolio.go           # Positions, Demat Holdings & Margin conversions
├── funds.go               # Account funds, margin calculations & limit checks
├── charts.go              # Historical Intraday & Daily OHLCV candle charts
├── data.go                # Market depth, quote snapshots & LTP
├── security.go            # Security master lookup & scrip details
├── instruments.go         # Bulk instrument contract master downloader
├── ledger.go              # Historical trade ledger & financial transactions
├── edis.go                # CDSL eDIS authorization & form generation
├── trades.go              # Trade book queries
├── user.go                # Client profile, KYC & segment configurations
├── postback.go            # Postback webhook parsing & validation
│
├── models/                # Pure domain data models & custom serialization
│   ├── orders.go          # OrderRequest, OrderModifyRequest, OrderResponse, OrderItem
│   ├── forever_orders.go  # ForeverOrderRequest, ForeverOrderResponse
│   ├── super_orders.go    # SuperOrderRequest, SuperOrderResponse
│   ├── portfolio.go       # Position, Holding, ConvertPositionRequest
│   ├── funds.go           # FundLimits, MarginCalculationRequest, MarginData
│   ├── charts.go          # HistoricalData, IntradayData, ToCandles() converter
│   ├── data.go            # Quote, MarketDepth, DepthItem
│   ├── security.go        # SecurityDetails, SecurityID
│   ├── instruments.go     # Instrument contract definitions
│   ├── ledger.go          # LedgerEntry structures
│   ├── edis.go            # EDISInquiryResponse, EDISFormRequest
│   ├── trades.go          # TradeItem structures
│   ├── user.go            # UserProfile, SegmentStatus
│   ├── postback.go        # OrderPostback structures
│   └── time.go            # Custom timestamp deserializer (models.Time)
│
├── ticker/                # Single-Socket Multiplexed WebSocket Streaming
│   ├── ticker.go          # Single-socket client, auto-reconnect & callback router
│   ├── parser.go          # Binary tick demuxer (LTP, Quote, Depth) & text order parser
│   └── types.go           # Tick, OrderUpdate, TickerSubscription
│
└── example/
    └── main.go            # End-to-end runnable demonstration
```

---

## Key Features

- **Standardized Namespace**: Packaged as `github.com/DragoHex/godhan`.
- **Zero Loss Financial Precision**: Strict use of `shopspring/decimal.Decimal` for prices, trigger levels, and account balances.
- **Single-Socket Multiplexed Ticker (`ticker`)**:
  - Connects to DhanHQ's unified binary WebSocket feed.
  - Demultiplexes binary quote packets (`OnTick`) and asynchronous JSON order postbacks (`OnOrderUpdate`) over a single persistent TCP connection.
- **Context First**: Full context propagation across all network requests.
- **Historical Chart Transformation**: Custom `models.HistoricalData.ToCandles()` method transforms columnar chart responses into standard time-series slices.
- **Advanced Order Varieties**: Full support for Forever (GTT), Bracket/Super orders, AMO, and Iceberg/Slicing.

---

## Installation

```sh
go get github.com/DragoHex/godhan
```

Requires Go `1.23` or later.

---

## Quickstart

### 1. Initialize Client & Query Balances

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DragoHex/godhan"
)

func main() {
	ctx := context.Background()

	// Initialize with your Dhan Client ID and Access Token
	client := godhan.NewClient("1000000001", "your-access-token")

	// Query Profile
	profile, err := client.GetUserProfile(ctx)
	if err != nil {
		log.Fatalf("GetUserProfile error: %v", err)
	}
	fmt.Printf("Logged in: Dhan Client ID = %s\n", profile.DhanClientID)

	// Query Available Funds
	funds, err := client.GetFundLimits(ctx)
	if err != nil {
		log.Fatalf("GetFundLimits error: %v", err)
	}
	fmt.Printf("Available Cash Balance: %s\n", funds.AvailableBalance.String())
}
```

---

### 2. Place & Manage Orders

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DragoHex/godhan"
	"github.com/DragoHex/godhan/models"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()
	client := godhan.NewClient("1000000001", "your-access-token")

	// Place Limit Buy Order
	orderResp, err := client.PlaceOrder(ctx, models.OrderRequest{
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		ProductType:     "INTRADAY",
		OrderType:       "LIMIT",
		Validity:        "DAY",
		TradingSymbol:   "TATASTEEL",
		SecurityID:      "11536",
		Quantity:        10,
		Price:           decimal.NewFromFloat(120.50),
	})
	if err != nil {
		log.Fatalf("PlaceOrder error: %v", err)
	}

	fmt.Printf("Order Submitted! OrderID: %s, Status: %s\n",
		orderResp.OrderID, orderResp.OrderStatus)

	// Query Positions
	positions, err := client.GetPositions(ctx)
	if err != nil {
		log.Fatalf("GetPositions error: %v", err)
	}
	for _, p := range positions {
		fmt.Printf("Position %s: NetQty = %d, Unrealized PnL = %s\n",
			p.TradingSymbol, p.NetQty, p.UnrealizedProfitLoss.String())
	}
}
```

---

### 3. Multiplexed Streaming Ticker (`ticker`)

Connect to DhanHQ's unified WebSocket feed for live market data and instant order alerts:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DragoHex/godhan/ticker"
)

func main() {
	// Initialize single-socket multiplexed ticker
	tk := ticker.NewTicker("1000000001", "your-access-token")

	tk.OnConnect(func() {
		fmt.Println("[WebSocket] Connected to DhanHQ Feed Gateway!")
		_ = tk.Subscribe([]ticker.TickerSubscription{
			{ExchangeSegment: 1, SecurityID: 11536}, // NSE EQ - TATASTEEL
		})
	})

	// Market Quote Packet Callback
	tk.OnTick(func(tick ticker.Tick) {
		fmt.Printf("[Tick] Security: %d, LTP: %s, Volume: %d\n",
			tick.SecurityID, tick.LTP.String(), tick.Volume)
	})

	// Real-Time Order Execution Callback
	tk.OnOrderUpdate(func(update ticker.OrderUpdate) {
		fmt.Printf("[Order Alert] %s -> Status: %s, FilledQty: %d\n",
			update.OrderID, update.OrderStatus, update.FilledQty)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = tk.ServeWithContext(ctx)
}
```

---

## Error Handling & Sentinels

```go
if err != nil {
	switch {
	case errors.Is(err, godhan.ErrUnauthorized):
		// Dhan token expired
	case errors.Is(err, godhan.ErrRateLimit):
		// Dhan rate limit hit
	case errors.Is(err, godhan.ErrNotFound):
		// Resource not found
	default:
		var apiErr *godhan.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("Dhan API Error [%s]: %s (HTTP %d)\n",
				apiErr.ErrorCode, apiErr.Message, apiErr.HTTPStatus)
		}
	}
}
```

---

## Testing & Verification

Run tests with race detection:

```sh
go test -race ./...
```

Run linter / compiler diagnostics:

```sh
go vet ./...
```

Build the runnable example:

```sh
go build -o /dev/null ./example/main.go
```

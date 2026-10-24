# Phase 5: Live Market Streaming & Binary Ticker

Part of Gokiteconnect-parity Go SDK for DhanHQ (`godhan`).

---

## 1. Scope
- **WebSocket Streaming Client**:
  - `Ticker` struct wrapping persistent WebSocket connection (`wss://api-feed.dhan.co`).
  - Thread-safe lifecycle: `Connect()`, `Close()`, `IsConnected()`.
  - Reconnection engine with exponential backoff (1s -> 30s) and random jitter.
  - Heartbeat / ping-pong watchdog to detect dead TCP sockets.
- **Dhan Binary Protocol Specification**:
  - Binary Request Packet:
    - Code 15: Subscribe / Unsubscribe LTP (Ticker mode)
    - Code 16: Subscribe / Unsubscribe Quote (Detailed mode)
    - Code 17: Subscribe / Unsubscribe Full Depth (Level 2/20 Depth)
    - Code 21: Server Disconnect / Authorization handshake
  - Binary Response Packet Decoders:
    - **Header**: 8-byte standard header (Feed Response Code, Message Length, Exchange Segment, Security ID).
    - **Feed Code 11 (LTP / Ticker)**: 8 bytes payload (LTP, Last Traded Time).
    - **Feed Code 12 (Quote / Market Data)**: 50 bytes payload (LTP, LTQ, LTT, Avg Traded Price, Volume, Total Buy Qty, Total Sell Qty, Open, High, Low, Close, Open Interest).
    - **Feed Code 13 (Full Depth / 5-depth Bids & Asks)**: Full order book snapshot with 5 bid and 5 ask levels (Price, Quantity, Orders count).
    - **Feed Code 14 (Previous Close / OI Data)**.
- **Event Listeners & Callbacks**:
  - `OnConnect(func())`
  - `OnTick(func(tick Tick))`
  - `OnError(func(err error))`
  - `OnClose(func(code int, reason string))`
  - `OnReconnect(func(attempt int, delay time.Duration))`

---

## 2. Test-Driven Development (TDD) Test Matrix

Tests MUST be written and fail first (`ticker_test.go`, `ticker_binary_test.go`):

### A. Binary Packet Parser (`ticker_binary_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Parse Code 11 (LTP packet) binary fixture | Unpacks exact decimal LTP price and uint32 timestamp without bit drift. |
| **Positive** | Parse Code 12 (Quote packet) binary fixture | Unpacks OHLC, volume, open interest, and turnover into typed `QuoteTick`. |
| **Positive** | Parse Code 13 (Full Depth packet) binary fixture | Unpacks 5 bid levels and 5 ask levels into structured depth arrays. |
| **Negative** | Truncated binary packet (< expected byte length) | Returns descriptive unpacking error without panic or buffer overrun. |
| **Negative** | Unknown / unsupported feed response code | Triggers `OnError` with unrecognized packet warning and safely skips bytes. |
| **Edge Case** | Little-endian vs Big-endian integer decoding | Tests against endianness spec of DhanHQ protocol (little-endian for numerical prices/quantities). |

### B. Connection Lifecycle & Event Loop (`ticker_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Successful handshake and authentication | Sends auth request on connect, triggers `OnConnect` callback. |
| **Positive** | Dynamic instrument subscription | Subscribes to instruments mid-session; verifies binary outgoing frame. |
| **Positive** | Clean shutdown | Calls `Close()`, terminates goroutines, verifies zero leaks via `goleak`. |
| **Negative** | Server drops connection unexpectedly | Auto-reconnect triggers with backoff; resumes active subscriptions. |
| **Negative** | Authentication token rejected by streaming gateway | Emits `ErrAuthExpired` via `OnError` callback; stops reconnect loop. |
| **Edge Case** | Rapid consecutive Subscribe / Unsubscribe calls | Channel buffers and mutex locks prevent data races (`go test -race`). |

# Master TDD Test Matrix: DhanHQ Go SDK (`godhan`)

This document defines the quality gate and mandatory Test-Driven Development (TDD) workflow for all phases of `godhan`.

---

## 1. TDD Workflow: Red-Green-Refactor Rule
1. **Red Stage**: Write unit tests for the target functionality covering positive, negative, and edge cases *before* any implementation code is written. Verify tests fail compilation or execution (`go test ./...` fails).
2. **Green Stage**: Implement minimum required production code to satisfy the tests. Verify tests pass cleanly (`go test ./...` succeeds).
3. **Refactor Stage**: Optimize memory allocations, clean up abstractions, check docstrings and comments. Run `go test -race ./...` and `go vet ./...` to guarantee thread safety and clean linter status.

---

## 2. Test Suites by Phase

| Phase | Test Suite | Target Areas | Key Verification Points |
| :--- | :--- | :--- | :--- |
| **Phase 1** | `errors_test.go`<br>`transport_test.go`<br>`time_test.go` | Client, Transport, Headers, Errors, Timestamps | Sentinels (`ErrInsufficientFunds`, `ErrOrderNotModifiable`, etc.), HTTP 200 failure detection, ISO8601 & space-separated dates. |
| **Phase 2** | `security_test.go`<br>`edis_test.go` | IP Setup, Kill Switch, PnL Exit, EDIS TPIN & Forms | Pre-flight input validations (IPv4 format, negative thresholds), idempotency of kill switch, TPIN generation. |
| **Phase 3** | `orders_test.go`<br>`super_orders_test.go`<br>`forever_orders_test.go`<br>`trades_test.go` | Regular orders, Slicing, Super (bracket), Forever (GTT), Trades | Quantity/price validations, order leg cancellation, freeze limit slicing, trade book pagination. |
| **Phase 4** | `portfolio_test.go`<br>`funds_test.go`<br>`charts_test.go`<br>`data_test.go` | Holdings, Positions, Funds, Margins, Daily/Intraday Candles | Safe zero-slice deserialization, decimal precision preservation, single & multi-leg margin computation, OHLC intervals. |
| **Phase 5** | `ticker_binary_test.go`<br>`ticker_test.go` | Binary packet parsing (Codes 11, 12, 13), WebSocket loop, Watchdog | Byte-level binary unpacking, little-endian numerical integrity, auto-reconnect backoff, goroutine leak checks (`goleak`). |

---

## 3. Strict Quality Gates
- **Host Go Version**: Latest installed Go on system (`go1.27.1 darwin/arm64`).
- **Race Condition Immunity**: All tests must execute cleanly under `go test -race ./...`.
- **Zero Linters / Vet Errors**: `go vet ./...` must return 0 warnings/errors.
- **Precision**: Zero `float64` for money, price, fees, and lot sizes (`shopspring/decimal.Decimal` required).
- **IP Protection**: Clean-room implementation; zero Zerodha/Kite IP terminology or copy.

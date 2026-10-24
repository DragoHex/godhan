# Phase 4: Portfolio, Funds, Statements, Charts & Analytics

Part of Gokiteconnect-parity Go SDK for DhanHQ (`godhan`).

---

## 1. Scope
- **Portfolio & Positions**:
  - `GetHoldings(ctx)` -> `GET /holdings` (Demat equity holdings, collateral, T1 holdings)
  - `GetPositions(ctx)` -> `GET /positions` (Day/net positions, open lots, realized/unrealized PnL)
  - `ConvertPosition(ctx, req)` -> `POST /positions/convert` (CNC <-> MIS, NRML <-> MIS)
  - `ExitAllPositions(ctx)` -> `DELETE /positions` (Square off all open intraday/margin positions)
- **Funds & Margin Calculators**:
  - `GetFundLimits(ctx)` -> `GET /fundlimit` (Available cash, collateral, margin utilized, payin/payout)
  - `CalculateMargin(ctx, req)` -> `POST /margincalculator` (Single order margin requirements)
  - `CalculateMultiMargin(ctx, req)` -> `POST /margincalculator/multi` (Basket / multi-leg margin requirements)
- **Statements & Ledger**:
  - `GetLedger(ctx, fromDate, toDate)` -> `GET /ledger` (Financial transaction statement)
- **Market Data, Charts & Analytics**:
  - `GetHistoricalCharts(ctx, req)` -> `POST /charts/historical` (Daily OHLC candle arrays)
  - `GetIntradayCharts(ctx, req)` -> `POST /charts/intraday` (Minute OHLC candle arrays)
  - `GetRollingOptionCharts(ctx, req)` -> `POST /charts/rollingoption` (Historical expired contract charts)
  - `GetCompanyInfo(ctx, req)` -> `POST /data/companyinfo`
  - `GetMarketMovers(ctx, req)` -> `POST /data/marketmovers` (Gainers, Losers, 52W High/Low)
  - `GetNewsHeadlines(ctx, req)` -> `POST /data/newsheadline`
  - `GetTechnicalData(ctx, req)` -> `POST /data/technical` (Pivot points, SMA, EMA indicators)
  - `GetInstruments(ctx)` -> Streaming CSV download and parser for Dhan daily scrip master (`https://images.dhan.co/api-data/api-scrip-master.csv`), adhering to KiteSDK streaming CSV design principles.

---

## 2. Test-Driven Development (TDD) Test Matrix

Tests MUST be written and fail first (`portfolio_test.go`, `funds_test.go`, `charts_test.go`, `data_test.go`):

### A. Portfolio & Positions (`portfolio_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Fetch demat holdings | Parses holdings into `[]Holding` with decimal quantities, average prices, and collateral flags. |
| **Positive** | Fetch open positions | Parses open net and day positions with mark-to-market and realized/unrealized PnL. |
| **Positive** | Convert position CNC -> MIS | Sends conversion payload, returns updated position status. |
| **Positive** | Exit all open positions | Sends square-off command, returns closed order responses. |
| **Negative** | Convert position with quantity greater than open net position | Returns `*APIError` with `ErrInvalidInput`. |
| **Edge Case** | Zero holdings / zero positions | Safely deserializes to empty slice `[]Holding{}` / `[]Position{}` without nil-pointer panics. |

### B. Funds & Margin Calculators (`funds_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Fetch fund limit details | Returns `FundLimit` struct with decimal cash balance, margin utilized, and collateral. |
| **Positive** | Compute single order margin | Returns required margin, span margin, and exposure margin. |
| **Positive** | Compute multi-leg basket margin | Calculates portfolio hedge benefit and final combined margin requirement. |
| **Negative** | Calculate margin with invalid exchange or security ID | Returns `ErrInvalidInput`. |
| **Edge Case** | High-precision decimal cash values (e.g., fractional paise) | `shopspring/decimal.Decimal` preserves exact numerical precision without IEEE-754 float drift. |

### C. Charts & Market Analytics (`charts_test.go`, `data_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Historical daily OHLC data | Deserializes timestamps, open, high, low, close, volume into `[]Candle`. |
| **Positive** | Intraday 1-minute / 5-minute candle data | Parses timestamped intervals into structured candles. |
| **Positive** | Technical indicators query | Parses SMA, EMA, and pivot points. |
| **Positive** | Market movers query | Returns top gainers and losers for index/exchange. |
| **Negative** | Intraday chart request exceeding allowed date window | Returns `*APIError` with descriptive message. |
| **Edge Case** | Holiday or non-trading hours empty candle array | Returns empty slice without parsing failure. |

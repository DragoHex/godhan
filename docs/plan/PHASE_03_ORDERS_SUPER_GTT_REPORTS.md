# Phase 3: Regular Orders, Super (Bracket), Forever (GTT), Slicing & Trade Reports

Part of Gokiteconnect-parity Go SDK for DhanHQ (`godhan`).

---

## 1. Scope
- **Regular Order Lifecycle**:
  - `PlaceOrder(ctx, req)` -> `POST /orders`
  - `ModifyOrder(ctx, orderID, req)` -> `PUT /orders/{order-id}`
  - `CancelOrder(ctx, orderID)` -> `DELETE /orders/{order-id}`
  - `SliceOrder(ctx, req)` -> `POST /orders/slicing` (auto-split large quantities per exchange freeze limits)
  - `GetOrders(ctx)` -> `GET /orders`
  - `GetOrderByID(ctx, orderID)` -> `GET /orders/{order-id}`
  - `GetOrderByCorrelationID(ctx, correlationID)` -> `GET /orders/external/{correlation-id}`
- **Super Orders (Bracket Orders with Target & StopLoss legs)**:
  - `PlaceSuperOrder(ctx, req)` -> `POST /super/orders`
  - `ModifySuperOrder(ctx, orderID, req)` -> `PUT /super/orders/{order-id}`
  - `CancelSuperOrder(ctx, orderID, leg)` -> `DELETE /super/orders/{order-id}/{order-leg}`
  - `GetSuperOrders(ctx)` -> `GET /super/orders`
- **Forever Orders (Good-Til-Triggered / GTT)**:
  - `PlaceForeverOrder(ctx, req)` -> `POST /forever/orders`
  - `ModifyForeverOrder(ctx, orderID, req)` -> `PUT /forever/orders/{order-id}`
  - `CancelForeverOrder(ctx, orderID)` -> `DELETE /forever/orders/{order-id}`
  - `GetForeverOrders(ctx)` -> `GET /forever/orders`
- **Conditional & Multi Orders (Alert Orders)**:
  - `PlaceAlertOrder(ctx, req)` -> `POST /alerts/orders`
  - `GetAlertOrders(ctx)` -> `GET /alerts/orders`
  - `GetAlertOrderByID(ctx, alertID)` -> `GET /alerts/orders/{alertId}`
  - `ModifyAlertOrder(ctx, alertID, req)` -> `PUT /alerts/orders/{alertId}`
  - `CancelAlertOrder(ctx, alertID)` -> `DELETE /alerts/orders/{alertId}`
  - `PlaceMultiOrder(ctx, req)` -> `POST /alerts/multi/orders`
- **Trade Book & Audit Trail**:
  - `GetTrades(ctx)` -> `GET /trades`
  - `GetTradesByOrderID(ctx, orderID)` -> `GET /trades/{order-id}`
  - `GetTradeHistory(ctx, fromDate, toDate, page)` -> `GET /trades/{from-date}/{to-date}/{page-number}`

---

## 2. Test-Driven Development (TDD) Test Matrix

Tests MUST be written and fail first (`orders_test.go`, `super_orders_test.go`, `forever_orders_test.go`, `trades_test.go`):

### A. Regular Orders & Slicing (`orders_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Place standard LIMIT BUY order | Returns `OrderResponse` with `orderId` and status `PENDING`/`TRANSIT`. |
| **Positive** | Slice order above freeze limit (e.g., 5000 qty NIFTY) | Calls `/orders/slicing`, returns array of child order IDs. |
| **Positive** | Modify pending order price and trigger price | Sends `PUT /orders/{id}`, returns updated order payload. |
| **Positive** | Cancel pending order | Sends `DELETE /orders/{id}`, returns cancelled confirmation. |
| **Positive** | Query order by external correlation ID | Returns matched order without needing internal broker order ID. |
| **Negative** | Place order with zero/negative quantity or price | Pre-flight validation fails with `ErrInvalidInput`. |
| **Negative** | Modify already completed/cancelled order | Returns `*APIError` wrapping `ErrOrderNotModifiable` (DH-905). |
| **Negative** | Insufficient margin/funds response | Returns `*APIError` wrapping `ErrInsufficientFunds` (DH-901). |
| **Edge Case** | AMO (After Market Order) placement outside trading hours | Flags AMO parameter accurately in request payload. |

### B. Super Orders (Bracket Orders) (`super_orders_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Place bracket order with entry, stop-loss and target | Sends `POST /super/orders`, returns parent and leg references. |
| **Positive** | Modify target price on active super order | Sends `PUT /super/orders/{id}` with leg identifier. |
| **Positive** | Cancel stop-loss leg or entire bracket order | Sends `DELETE /super/orders/{id}/{leg}`, verifies leg cancellation. |
| **Negative** | Place super order where target <= entry on BUY | Pre-flight validation returns `ErrInvalidInput`. |
| **Negative** | Cancel non-existent super order leg | Returns `*APIError` wrapping `ErrNotFound`. |
| **Edge Case** | Trailing stop-loss step value verification | Verifies decimal step resolution against instrument tick size. |

### C. Forever Orders (GTT) & Alerts (`forever_orders_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Place single-condition GTT forever order | Sends `POST /forever/orders`, returns GTT order identifier. |
| **Positive** | Place OCO (One-Cancels-Other) two-condition GTT | Sends upper and lower price triggers with respective order legs. |
| **Positive** | Fetch all active forever orders | Returns slice of `ForeverOrder` models. |
| **Positive** | Cancel active forever order | Sends `DELETE /forever/orders/{id}`, returns confirmation. |
| **Negative** | GTT trigger price inverted or invalid | Validation triggers `ErrInvalidInput`. |
| **Edge Case** | Empty order book returns empty slice `[]Order{}` | No nil slice serialization surprises. |

### D. Trade Book & History (`trades_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Fetch current day trades (`GET /trades`) | Returns list of executed trades with fills, execution prices, exchange trade IDs. |
| **Positive** | Fetch trades for specific order ID | Returns list of fills for that order. |
| **Positive** | Fetch paginated trade history | Requests `{from}/{to}/{page}`, parses items and pagination cursor. |
| **Negative** | Historical trade request with invalid date range | Pre-flight validates `from <= to` or returns `ErrInvalidInput`. |

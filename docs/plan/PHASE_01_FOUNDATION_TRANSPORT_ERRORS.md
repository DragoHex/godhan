# Phase 1: Foundation, Transport, Models & Error Taxonomy

Part of Gokiteconnect-parity Go SDK for DhanHQ (`godhan`).

---

## 1. Scope
- Base `Client` struct with configuration options (`WithBaseURI`, `WithHTTPClient`, `WithAccessToken`, `WithClientID`).
- Transport seam (`HTTPClient` interface) with keep-alive connection pooling (`MaxIdleConnsPerHost: 50`, `ResponseHeaderTimeout: 5s`).
- Header injection: `access-token`, `client-id`, `Content-Type: application/json`, `Accept: application/json`.
- Modern Go 1.13+ Error taxonomy (`APIError`, `ErrAuthExpired`, `ErrInsufficientFunds`, `ErrOrderNotModifiable`, `ErrInvalidInput`, `ErrRateLimited`, `ErrBrokerUnavailable`).
- Strongly-typed models with `shopspring/decimal.Decimal` and custom timestamp unmarshaler (`BrokerTime`).
- Package URI constants in `endpoints.go`.

---

## 2. Test-Driven Development (TDD) Test Matrix

Tests MUST be written and fail first (`errors_test.go`, `transport_test.go`, `time_test.go`):

### A. Error Taxonomy (`errors_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Dhan returns error JSON with `status: "failure"`, `error_code: "DH-901"`, `message: "Insufficient Funds"` | Returns `*APIError` with `errors.Is(err, ErrInsufficientFunds) == true`. |
| **Positive** | Dhan returns `errorType: "OrderError"`, `errorCode: "DH-905"`, `errorMessage: "Order is already closed"` | Returns `*APIError` with `errors.Is(err, ErrOrderNotModifiable) == true`. |
| **Positive** | Gateway returns HTTP 401 Unauthorized with token expired message | Returns `*APIError` with `errors.Is(err, ErrAuthExpired) == true`. |
| **Negative** | Gateway returns HTTP 502 Bad Gateway with raw HTML string body | Returns synthetic `*APIError` wrapping `ErrBrokerUnavailable`. |
| **Negative** | Dhan returns HTTP 200 OK but payload contains `status: "failure"` | Correctly identifies failure despite HTTP 200 status code. |
| **Edge Case** | Malformed JSON or truncated error body | Returns fallback error without panic or nil-pointer dereference. |

### B. Transport & Headers (`transport_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Request headers injection | Injects `access-token`, `client-id`, `Content-Type: application/json`, `Accept: application/json`. |
| **Positive** | Request JSON serialization | Encodes body payload cleanly to JSON bytes. |
| **Negative** | HTTP client returns network timeout | Wraps and returns `ErrBrokerUnavailable`. |
| **Edge Case** | Query parameters with special characters | Properly encodes query string in GET/DELETE requests. |

### C. Timestamp Deserialization (`time_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Standard Dhan datetime `"2023-10-15 14:30:00"` | Parsed into exact `time.Time`. |
| **Positive** | ISO8601 `"2023-10-15T14:30:00+05:30"` | Parsed into exact `time.Time`. |
| **Positive** | Date only `"2023-10-15"` | Parsed into midnight `time.Time`. |
| **Edge Case** | `null` or empty string `""` | Sets zero time without unmarshaling error. |
| **Negative** | Corrupt date string | Returns descriptive error. |

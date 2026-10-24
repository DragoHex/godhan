# Phase 2: Auth, Security, Trader Controls & EDIS

Part of Gokiteconnect-parity Go SDK for DhanHQ (`godhan`).

---

## 1. Scope
- **IP Whitelisting Management**:
  - `GetIP(ctx)` -> `GET /ip/getIP`
  - `SetIP(ctx, req)` -> `POST /ip/setIP`
  - `ModifyIP(ctx, req)` -> `PUT /ip/modifyIP`
- **Trader Controls**:
  - **Kill Switch**:
    - `GetKillSwitchStatus(ctx)` -> `GET /killswitch`
    - `SetKillSwitch(ctx, status)` -> `POST /killswitch`
  - **PnL-based Exit**:
    - `GetPnLExit(ctx)` -> `GET /pnlExit`
    - `SetPnLExit(ctx, req)` -> `POST /pnlExit`
    - `StopPnLExit(ctx)` -> `DELETE /pnlExit`
- **Electronic Delivery Instruction Slip (EDIS)**:
  - `GenerateTPIN(ctx)` -> `GET /edis/tpin`
  - `GenerateEDISForm(ctx, req)` -> `POST /edis/form`
  - `GenerateBulkEDISForm(ctx, req)` -> `POST /edis/bulkform`
  - `InquireEDISQuantity(ctx, isin)` -> `GET /edis/inquire/{isin}`

---

## 2. Test-Driven Development (TDD) Test Matrix

Tests MUST be written and fail first (`security_test.go`, `edis_test.go`):

### A. IP Setup & Whitelisting (`security_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Fetch registered static IP (`GET /ip/getIP`) | Returns `IPDetails` struct with parsed IP list and status. |
| **Positive** | Set new static IP (`POST /ip/setIP`) | Sends `{ "ip": "203.0.113.1" }`, returns success status. |
| **Positive** | Modify existing static IP (`PUT /ip/modifyIP`) | Sends modified IP payload, returns confirmation. |
| **Negative** | Invalid IPv4/IPv6 address format | Client-side validation returns `ErrInvalidInput` before network call. |
| **Negative** | Token expired during IP modification | Returns `*APIError` wrapping `ErrAuthExpired`. |
| **Edge Case** | Multiple comma-separated IPs or CIDR notation | Validates against Dhan schema rules. |

### B. Trader Controls: Kill Switch & PnL Exit (`security_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Check Kill Switch status (`GET /killswitch`) | Returns `KillSwitchStatus` (`ACTIVATED` / `DEACTIVATED`). |
| **Positive** | Activate Kill Switch (`POST /killswitch`) | Sends activation flag; locks account from placing new orders. |
| **Positive** | Configure PnL exit thresholds (`POST /pnlExit`) | Sends profit and loss triggers, returns configured parameters. |
| **Positive** | Stop PnL exit monitor (`DELETE /pnlExit`) | Clears automatic square-off rule. |
| **Negative** | Setting negative profit target or inverted loss limit | Returns `ErrInvalidInput`. |
| **Edge Case** | Dhan returns 200 with remarks `"Kill switch already active"` | Idempotently returns current state without error. |

### C. EDIS Authorization & TPIN (`edis_test.go`)
| Case Type | Scenario | Expected Behavior |
| :--- | :--- | :--- |
| **Positive** | Request TPIN generation (`GET /edis/tpin`) | Triggers SMS/Email OTP, returns success flag. |
| **Positive** | Inquire authorized quantity for ISIN (`GET /edis/inquire/INE002A01018`) | Returns `EDISStatus` with authorized vs allocated quantities. |
| **Positive** | Generate HTML EDIS Form (`POST /edis/form`) | Returns CDSL redirect HTML form string and target URL. |
| **Positive** | Generate Bulk EDIS Form (`POST /edis/bulkform`) | Accepts multiple ISINs, returns batch form payload. |
| **Negative** | Query EDIS status with invalid/empty ISIN | Returns `ErrInvalidInput`. |
| **Negative** | TPIN generation rate-limited by CDSL | Returns `*APIError` wrapping `ErrRateLimited`. |

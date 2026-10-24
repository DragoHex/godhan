package godhan

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Standard Sentinel Domain Errors
var (
	ErrAuthExpired        = errors.New("authentication expired or invalid")
	ErrInsufficientFunds  = errors.New("insufficient funds in trading account")
	ErrOrderNotModifiable = errors.New("order is closed or cannot be modified")
	ErrInvalidInput       = errors.New("invalid request parameters")
	ErrRateLimited        = errors.New("rate limit exceeded")
	ErrBrokerUnavailable  = errors.New("broker gateway unavailable or timed out")
	ErrOrderRejected      = errors.New("order rejected by broker RMS or exchange")
	ErrNotFound           = errors.New("resource not found")
	ErrPermissionDenied   = errors.New("permission denied for exchange or product")
)

// APIError represents structured error from DhanHQ gateway
type APIError struct {
	HTTPStatus int
	ErrorCode  string
	ErrorType  string
	Message    string
	RawBody    []byte
	Err        error
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("dhan api error [%s] (HTTP %d): %s", e.ErrorCode, e.HTTPStatus, e.Message)
	}
	return fmt.Sprintf("dhan api error (HTTP %d): %s", e.HTTPStatus, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// rawDhanError represents potential error response envelopes from DhanHQ
type rawDhanError struct {
	Status       string          `json:"status"`
	Remarks      json.RawMessage `json:"remarks"`
	ErrorType    string          `json:"errorType"`
	ErrorCode    string          `json:"errorCode"`
	ErrorMessage string          `json:"errorMessage"`
	Data         interface{}     `json:"data"`
}

type rawRemarksObj struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// parseErrorResponse unmarshals DhanHQ error payloads and maps them to domain sentinels
func parseErrorResponse(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		HTTPStatus: statusCode,
		RawBody:    body,
	}

	if len(body) == 0 {
		apiErr.Message = http.StatusText(statusCode)
		apiErr.Err = mapHTTPStatusToSentinel(statusCode)
		return apiErr
	}

	var raw rawDhanError
	if err := json.Unmarshal(body, &raw); err != nil {
		// Non-JSON response (e.g. 502 Bad Gateway HTML, proxy errors)
		apiErr.Message = string(body)
		apiErr.Err = mapHTTPStatusToSentinel(statusCode)
		return apiErr
	}

	apiErr.ErrorType = raw.ErrorType
	apiErr.ErrorCode = raw.ErrorCode
	apiErr.Message = raw.ErrorMessage

	// Parse remarks field which can be either an object {"error_code":"...", "message":"..."} or a raw string
	if len(raw.Remarks) > 0 {
		var remarksObj rawRemarksObj
		if err := json.Unmarshal(raw.Remarks, &remarksObj); err == nil && remarksObj.Message != "" {
			if apiErr.ErrorCode == "" {
				apiErr.ErrorCode = remarksObj.ErrorCode
			}
			if apiErr.Message == "" {
				apiErr.Message = remarksObj.Message
			}
		} else {
			var remarksStr string
			if err := json.Unmarshal(raw.Remarks, &remarksStr); err == nil {
				if apiErr.Message == "" {
					apiErr.Message = remarksStr
				}
			}
		}
	}

	if apiErr.Message == "" && raw.Status != "" {
		apiErr.Message = fmt.Sprintf("status: %s", raw.Status)
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(statusCode)
	}

	apiErr.Err = classifyError(statusCode, apiErr.ErrorCode, apiErr.ErrorType, apiErr.Message)
	return apiErr
}

func mapHTTPStatusToSentinel(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrAuthExpired
	case http.StatusForbidden:
		return ErrPermissionDenied
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusBadRequest:
		return ErrInvalidInput
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return ErrBrokerUnavailable
	default:
		if statusCode >= 500 {
			return ErrBrokerUnavailable
		}
		return ErrOrderRejected
	}
}

func classifyError(statusCode int, code, errType, msg string) error {
	lowerMsg := strings.ToLower(msg)
	lowerCode := strings.ToUpper(code)

	// Explicit Dhan error code mappings
	if lowerCode == "DH-901" || strings.Contains(lowerMsg, "insufficient funds") || strings.Contains(lowerMsg, "insufficient balance") {
		return ErrInsufficientFunds
	}
	if lowerCode == "DH-905" || strings.Contains(lowerMsg, "order is already closed") || strings.Contains(lowerMsg, "cannot be modified") {
		return ErrOrderNotModifiable
	}
	if statusCode == http.StatusUnauthorized || strings.Contains(lowerMsg, "token expired") || strings.Contains(lowerMsg, "invalid jwt") {
		return ErrAuthExpired
	}
	if statusCode == http.StatusTooManyRequests || strings.Contains(lowerMsg, "too many requests") || strings.Contains(lowerMsg, "rate limit") {
		return ErrRateLimited
	}
	if statusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if statusCode == http.StatusForbidden || strings.Contains(lowerMsg, "permission") || strings.Contains(lowerMsg, "not authorized") {
		return ErrPermissionDenied
	}
	if statusCode == http.StatusBadRequest && (errType == "InputException" || strings.Contains(lowerMsg, "invalid")) {
		return ErrInvalidInput
	}
	if statusCode >= 500 {
		return ErrBrokerUnavailable
	}

	return ErrOrderRejected
}

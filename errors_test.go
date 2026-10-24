package godhan

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorTaxonomy_Sentinels(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		targetErr  error
	}{
		{
			name:       "DH-901 Insufficient Funds in Remarks",
			statusCode: http.StatusBadRequest,
			body:       `{"status":"failure","remarks":{"error_code":"DH-901","message":"Insufficient Funds in Account"},"data":null}`,
			targetErr:  ErrInsufficientFunds,
		},
		{
			name:       "DH-905 models.Order Not Modifiable",
			statusCode: http.StatusBadRequest,
			body:       `{"errorType":"OrderError","errorCode":"DH-905","errorMessage":"models.Order is already closed or cancelled"}`,
			targetErr:  ErrOrderNotModifiable,
		},
		{
			name:       "HTTP 401 Unauthorized Auth Expired",
			statusCode: http.StatusUnauthorized,
			body:       `{"status":"failure","remarks":"Token expired or invalid JWT"}`,
			targetErr:  ErrAuthExpired,
		},
		{
			name:       "HTTP 400 Invalid Input",
			statusCode: http.StatusBadRequest,
			body:       `{"errorType":"InputException","errorCode":"DH-801","errorMessage":"Invalid price tick"}`,
			targetErr:  ErrInvalidInput,
		},
		{
			name:       "HTTP 429 Rate Limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"status":"failure","remarks":"Too Many Requests: rate limit exceeded"}`,
			targetErr:  ErrRateLimited,
		},
		{
			name:       "HTTP 502 Gateway Unavailable",
			statusCode: http.StatusBadGateway,
			body:       `<html><body>502 Bad Gateway</body></html>`,
			targetErr:  ErrBrokerUnavailable,
		},
		{
			name:       "HTTP 200 OK with status failure",
			statusCode: http.StatusOK,
			body:       `{"status":"failure","remarks":"models.Order placement rejected by RMS","data":null}`,
			targetErr:  ErrOrderRejected,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			apiErr := parseErrorResponse(tc.statusCode, []byte(tc.body))
			if apiErr == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(apiErr, tc.targetErr) {
				t.Errorf("expected errors.Is(err, %v) == true, got %v (wrapped: %v)", tc.targetErr, apiErr, apiErr.Unwrap())
			}
		})
	}
}

func TestErrorTaxonomy_EdgeCases(t *testing.T) {
	t.Run("Empty body", func(t *testing.T) {
		err := parseErrorResponse(http.StatusInternalServerError, []byte{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrBrokerUnavailable) {
			t.Errorf("expected ErrBrokerUnavailable, got %v", err)
		}
	})

	t.Run("Truncated malformed JSON", func(t *testing.T) {
		err := parseErrorResponse(http.StatusBadRequest, []byte(`{"status":"failure","remarks":{`))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() == "" {
			t.Errorf("expected non-empty error message")
		}
	})
}

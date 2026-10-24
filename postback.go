package godhan

import (
	"github.com/DragoHex/godhan/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

)

// models.PostbackPayload represents the webhook callback event pushed by DhanHQ upon order status transitions

// ParsePostback unmarshals raw JSON bytes received from Dhan webhook into models.PostbackPayload
func ParsePostback(body []byte) (*models.PostbackPayload, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("%w: empty postback body", ErrInvalidInput)
	}

	var payload models.PostbackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal postback JSON: %v", ErrInvalidInput, err)
	}

	return &payload, nil
}

// ParsePostbackRequest extracts and decodes models.PostbackPayload from incoming HTTP request
func ParsePostbackRequest(r *http.Request) (*models.PostbackPayload, error) {
	if r == nil || r.Body == nil {
		return nil, fmt.Errorf("%w: nil request or body", ErrInvalidInput)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read postback body: %v", ErrInvalidInput, err)
	}

	return ParsePostback(body)
}

// IsTerminalStatus checks if the given Dhan orderStatus string represents a final state
func IsTerminalStatus(status string) bool {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "TRADED", "REJECTED", "CANCELLED", "EXPIRED":
		return true
	default:
		return false
	}
}

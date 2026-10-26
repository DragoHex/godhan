package godhan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient interface abstracts the underlying HTTP transport seam
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// defaultHTTPClient constructs a high-performance HTTP client with connection pooling
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
}

// genericEnvelope captures standard Dhan envelope wrapper if present
type genericEnvelope struct {
	Status  string          `json:"status"`
	Remarks json.RawMessage `json:"remarks"`
	Data    json.RawMessage `json:"data"`
}

// doRequest performs an HTTP request, injects required Dhan headers, and parses the response
func (c *Client) doRequest(ctx context.Context, method, endpoint string, query url.Values, body interface{}, target interface{}) error {
	fullURL := strings.TrimRight(c.baseURI, "/") + "/" + strings.TrimLeft(endpoint, "/")
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		reqBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: failed to marshal request body: %v", ErrInvalidInput, err)
		}
		bodyReader = bytes.NewReader(reqBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %v", ErrInvalidInput, err)
	}

	clientID := c.ClientID()
	accessToken := c.AccessToken()

	// Inject mandatory Dhan headers
	req.Header.Set("client-id", clientID)
	req.Header.Set("dhanClientId", clientID)
	req.Header.Set("access-token", accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: network error: %v", ErrBrokerUnavailable, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: failed to read response body: %v", ErrBrokerUnavailable, err)
	}

	// 1. Check HTTP Status Code failure
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp.StatusCode, respBytes)
	}

	// 2. Check for Dhan HTTP 200 payload failure (e.g. status: "failure")
	var env genericEnvelope
	if err := json.Unmarshal(respBytes, &env); err == nil && strings.EqualFold(env.Status, "failure") {
		return parseErrorResponse(resp.StatusCode, respBytes)
	}

	// 3. Unpack into target object if target is provided
	if target != nil {
		// If envelope has "data" field populated, try unmarshaling "data" first
		if len(env.Data) > 0 && string(env.Data) != "null" {
			if err := json.Unmarshal(env.Data, target); err == nil {
				return nil
			}
		}
		// Direct unmarshal fallback
		if err := json.Unmarshal(respBytes, target); err != nil {
			return fmt.Errorf("failed to unmarshal response into target: %w", err)
		}
	}

	return nil
}

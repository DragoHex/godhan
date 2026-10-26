package godhan

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestTransport_HeaderInjection(t *testing.T) {
	var capturedHeader http.Header
	var capturedPath string
	var capturedMethod string
	var capturedBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Clone()
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		capturedBody = string(bodyBytes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","data":{"id":"12345"}}`))
	}))
	defer ts.Close()

	client := NewClient("mock-client-id", "mock-access-token",
		WithBaseURI(ts.URL),
	)

	ctx := context.Background()
	reqBody := map[string]string{"foo": "bar"}
	var result map[string]interface{}

	err := client.doRequest(ctx, http.MethodPost, "/test-path", nil, reqBody, &result)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if capturedHeader.Get("client-id") != "mock-client-id" {
		t.Errorf("expected client-id header 'mock-client-id', got '%s'", capturedHeader.Get("client-id"))
	}
	if capturedHeader.Get("dhanClientId") != "mock-client-id" {
		t.Errorf("expected dhanClientId header 'mock-client-id', got '%s'", capturedHeader.Get("dhanClientId"))
	}
	if capturedHeader.Get("access-token") != "mock-access-token" {
		t.Errorf("expected access-token header 'mock-access-token', got '%s'", capturedHeader.Get("access-token"))
	}
	if capturedHeader.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got '%s'", capturedHeader.Get("Content-Type"))
	}
	if capturedHeader.Get("Accept") != "application/json" {
		t.Errorf("expected Accept application/json, got '%s'", capturedHeader.Get("Accept"))
	}
	if capturedPath != "/test-path" {
		t.Errorf("expected path /test-path, got %s", capturedPath)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", capturedMethod)
	}
	if capturedBody != `{"foo":"bar"}` {
		t.Errorf("expected body {\"foo\":\"bar\"}, got %s", capturedBody)
	}
}

func TestTransport_QueryParams(t *testing.T) {
	var capturedQuery url.Values

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","data":[]}`))
	}))
	defer ts.Close()

	client := NewClient("mock-client", "mock-token", WithBaseURI(ts.URL))

	ctx := context.Background()
	params := url.Values{}
	params.Set("symbol", "RELIANCE")
	params.Set("limit", "10")

	var result []interface{}
	err := client.doRequest(ctx, http.MethodGet, "/query-test", params, nil, &result)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if capturedQuery.Get("symbol") != "RELIANCE" || capturedQuery.Get("limit") != "10" {
		t.Errorf("query params not forwarded properly: got %v", capturedQuery)
	}
}

func TestTransport_NetworkTimeout(t *testing.T) {
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()

	customHTTP := &http.Client{Timeout: 10 * time.Millisecond}
	client := NewClient("client", "token",
		WithBaseURI(slowServer.URL),
		WithHTTPClient(customHTTP),
	)

	var result map[string]interface{}
	err := client.doRequest(context.Background(), http.MethodGet, "/slow", nil, nil, &result)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

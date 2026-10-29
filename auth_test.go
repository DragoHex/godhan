package godhan

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestAuth_RenewToken_Success(t *testing.T) {
	var capturedHeader http.Header
	var capturedMethod string
	var capturedPath string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Clone()
		capturedMethod = r.Method
		capturedPath = r.URL.Path

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"dhanClientId": "1000000001",
			"dhanClientName": "Mohit Panwar",
			"dhanClientUcc": "ABC12345",
			"givenPowerOfAttorney": true,
			"accessToken": "new-refreshed-token-999",
			"expiryTime": "2026-10-27 15:30:00"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "old-active-token", WithBaseURI(ts.URL))

	resp, err := client.RenewToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedMethod != http.MethodGet {
		t.Errorf("method = %q, want %q", capturedMethod, http.MethodGet)
	}
	if capturedPath != URIRenewToken {
		t.Errorf("path = %q, want %q", capturedPath, URIRenewToken)
	}
	if capturedHeader.Get("access-token") != "old-active-token" {
		t.Errorf("access-token header = %q, want 'old-active-token'", capturedHeader.Get("access-token"))
	}
	if capturedHeader.Get("dhanClientId") != "1000000001" {
		t.Errorf("dhanClientId header = %q, want '1000000001'", capturedHeader.Get("dhanClientId"))
	}

	// Verify unmarshaled response fields
	if resp.DhanClientID != "1000000001" {
		t.Errorf("DhanClientID = %q, want 1000000001", resp.DhanClientID)
	}
	if resp.DhanClientName != "Mohit Panwar" {
		t.Errorf("DhanClientName = %q, want 'Mohit Panwar'", resp.DhanClientName)
	}
	if resp.DhanClientUCC != "ABC12345" {
		t.Errorf("DhanClientUCC = %q, want 'ABC12345'", resp.DhanClientUCC)
	}
	if !resp.GivenPowerOfAttorney {
		t.Errorf("GivenPowerOfAttorney = false, want true")
	}
	if resp.AccessToken != "new-refreshed-token-999" {
		t.Errorf("AccessToken = %q, want 'new-refreshed-token-999'", resp.AccessToken)
	}
	if resp.ExpiryTime.Time.IsZero() {
		t.Errorf("ExpiryTime parsed as zero time")
	} else if resp.ExpiryTime.Time.Year() != 2026 {
		t.Errorf("ExpiryTime year = %d, want 2026", resp.ExpiryTime.Time.Year())
	}

	// Verify in-place client state mutation
	if client.AccessToken() != "new-refreshed-token-999" {
		t.Errorf("client.AccessToken() = %q, want 'new-refreshed-token-999'", client.AccessToken())
	}
}

func TestAuth_RenewToken_TokenOverride(t *testing.T) {
	var capturedToken string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedToken = r.Header.Get("access-token")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"dhanClientId": "1000000001",
			"accessToken": "override-renewed-token-xyz",
			"expiryTime": "2026-10-27T08:00:00"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "default-client-token", WithBaseURI(ts.URL))

	resp, err := client.RenewToken(context.Background(), "explicit-override-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedToken != "explicit-override-token" {
		t.Errorf("captured access-token = %q, want 'explicit-override-token'", capturedToken)
	}
	if resp.AccessToken != "override-renewed-token-xyz" {
		t.Errorf("resp.AccessToken = %q, want 'override-renewed-token-xyz'", resp.AccessToken)
	}
	if client.AccessToken() != "override-renewed-token-xyz" {
		t.Errorf("client.AccessToken() = %q, want 'override-renewed-token-xyz'", client.AccessToken())
	}
}

func TestAuth_RenewAccessToken_ParityAlias(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"dhanClientId": "1000000001",
			"accessToken": "parity-alias-token-123",
			"expiryTime": "2026-10-27T10:00:00+05:30"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "initial-token", WithBaseURI(ts.URL))

	resp, err := client.RenewAccessToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AccessToken != "parity-alias-token-123" {
		t.Errorf("resp.AccessToken = %q, want 'parity-alias-token-123'", resp.AccessToken)
	}
	if client.AccessToken() != "parity-alias-token-123" {
		t.Errorf("client.AccessToken() = %q, want 'parity-alias-token-123'", client.AccessToken())
	}
}

func TestAuth_RenewToken_LiveGatewayResponseFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"createTime": "2026-10-29T15:08:58.164",
			"expiryTime": "2026-10-30T15:08:58.163",
			"token": "dhan-live-token-777"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "old-token", WithBaseURI(ts.URL))

	resp, err := client.RenewToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Token != "dhan-live-token-777" {
		t.Errorf("resp.Token = %q, want 'dhan-live-token-777'", resp.Token)
	}
	if resp.AccessToken != "dhan-live-token-777" {
		t.Errorf("resp.AccessToken = %q, want 'dhan-live-token-777'", resp.AccessToken)
	}
	if client.AccessToken() != "dhan-live-token-777" {
		t.Errorf("client.AccessToken() = %q, want 'dhan-live-token-777'", client.AccessToken())
	}
}

func TestAuth_RenewToken_EmptyToken(t *testing.T) {
	client := NewClient("1000000001", "")

	// 1. Without token in client and without override
	_, err := client.RenewToken(context.Background())
	if err == nil {
		t.Fatal("expected error when renewing empty token, got nil")
	}
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Whitespace-only override
	_, err = client.RenewToken(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error with whitespace override, got nil")
	}
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuth_RenewToken_ExpiredToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{
			"status": "failure",
			"remarks": {
				"error_code": "DH-903",
				"message": "Token expired or invalid"
			},
			"data": null
		}`))
	}))
	defer ts.Close()

	initialToken := "expired-token-abc"
	client := NewClient("1000000001", initialToken, WithBaseURI(ts.URL))

	resp, err := client.RenewToken(context.Background())
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response on error, got %+v", resp)
	}
	if !errors.Is(err, ErrAuthExpired) {
		t.Errorf("expected ErrAuthExpired, got %v", err)
	}

	// Verify token was not modified on failure
	if client.AccessToken() != initialToken {
		t.Errorf("client.AccessToken() = %q, want %q", client.AccessToken(), initialToken)
	}
}

func TestAuth_RenewToken_FailureEnvelope200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "failure",
			"remarks": "Invalid token signature",
			"data": null
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "bad-token", WithBaseURI(ts.URL))

	_, err := client.RenewToken(context.Background())
	if err == nil {
		t.Fatal("expected error for HTTP 200 failure envelope, got nil")
	}
}

func TestAuth_RenewToken_NetworkFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := ts.URL
	ts.Close() // Close immediately to trigger connection refusal

	client := NewClient("1000000001", "valid-token", WithBaseURI(url))

	_, err := client.RenewToken(context.Background())
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
	if !errors.Is(err, ErrBrokerUnavailable) {
		t.Errorf("expected ErrBrokerUnavailable, got %v", err)
	}
}

func TestAuth_RenewToken_ThreadSafety(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"dhanClientId": "1000000001",
			"accessToken": "concurrent-token-456",
			"expiryTime": "2026-10-27 12:00:00"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "seed-token", WithBaseURI(ts.URL))

	var wg sync.WaitGroup
	ctx := context.Background()

	// Launch readers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = client.AccessToken()
				_ = client.ClientID()
				time.Sleep(100 * time.Microsecond)
			}
		}()
	}

	// Launch writers (renewing token)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, _ = client.RenewToken(ctx)
			}
		}()
	}

	wg.Wait()
}

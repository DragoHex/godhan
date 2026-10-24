package godhan

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser_GetUserProfile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != URIGetIP {
			t.Errorf("path = %q, want %q", r.URL.Path, URIGetIP)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"dhanClientId": "1000000001",
			"ipFlag": true,
			"ip": "203.0.113.1"
		}`))
	}))
	defer ts.Close()

	client := NewClient("1000000001", "valid-token", WithBaseURI(ts.URL))
	profile, err := client.GetUserProfile(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if profile.DhanClientID != "1000000001" {
		t.Errorf("DhanClientID = %q, want 1000000001", profile.DhanClientID)
	}
	if !profile.IPFlag {
		t.Errorf("IPFlag = false, want true")
	}
}

func TestClient_GetLoginURL(t *testing.T) {
	client := New(WithClientID("1000000001"))
	urlStr := client.GetLoginURL("partner-xyz")
	if urlStr == "" {
		t.Fatal("expected non-empty login URL")
	}
}

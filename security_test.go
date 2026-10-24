package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestSecurity_IPManagement(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == URIGetIP:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"primaryIP": "203.0.113.1",
				"secondaryIP": "203.0.113.2",
				"detectedIP": "203.0.113.1",
				"ipMatchStatus": "PRIMARY_MATCH",
				"ordersAllowed": true
			}`))
		case r.Method == http.MethodPost && r.URL.Path == URISetIP:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"SUCCESS","message":"IP set successfully"}`))
		case r.Method == http.MethodPut && r.URL.Path == URIModifyIP:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"SUCCESS","message":"IP modified successfully"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GetIP
	ipDetails, err := client.GetIP(ctx)
	if err != nil {
		t.Fatalf("GetIP failed: %v", err)
	}
	if ipDetails.PrimaryIP != "203.0.113.1" || !ipDetails.OrdersAllowed {
		t.Errorf("unexpected IP details: %+v", ipDetails)
	}

	// 2. SetIP
	setResp, err := client.SetIP(ctx, models.UserIPRequest{
		IP:     "203.0.113.1",
		IPFlag: "PRIMARY",
	})
	if err != nil {
		t.Fatalf("SetIP failed: %v", err)
	}
	if setResp.Status != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", setResp.Status)
	}

	// 3. ModifyIP
	modResp, err := client.ModifyIP(ctx, models.UserIPRequest{
		IP:     "203.0.113.2",
		IPFlag: "SECONDARY",
	})
	if err != nil {
		t.Fatalf("ModifyIP failed: %v", err)
	}
	if modResp.Status != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", modResp.Status)
	}
}

func TestSecurity_TraderControls(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == URIKillSwitch:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"dhanClientId":"test-client","killSwitchStatus":"ACTIVATE"}`))
		case r.Method == http.MethodPost && r.URL.Path == URIKillSwitch:
			statusParam := r.URL.Query().Get("killSwitchStatus")
			if statusParam != "ACTIVATE" && statusParam != "DEACTIVATE" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"failure","remarks":"invalid status"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"dhanClientId":"test-client","killSwitchStatus":"` + statusParam + `"}`))
		case r.Method == http.MethodGet && r.URL.Path == URIPnLExit:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"pnlExitStatus": "ACTIVE",
				"profit": 5000.50,
				"loss": 2000.25,
				"enableKillSwitch": true,
				"productType": ["INTRADAY"]
			}`))
		case r.Method == http.MethodPost && r.URL.Path == URIPnLExit:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"pnlExitStatus":"ACTIVE","message":"PNL exit created"}`))
		case r.Method == http.MethodDelete && r.URL.Path == URIPnLExit:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"SUCCESS","message":"PNL exit stopped"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GetKillSwitchStatus
	ks, err := client.GetKillSwitchStatus(ctx)
	if err != nil {
		t.Fatalf("GetKillSwitchStatus failed: %v", err)
	}
	if ks.KillSwitchStatus != "ACTIVATE" {
		t.Errorf("expected ACTIVATE, got %s", ks.KillSwitchStatus)
	}

	// 2. SetKillSwitch
	ksResp, err := client.SetKillSwitch(ctx, "DEACTIVATE")
	if err != nil {
		t.Fatalf("SetKillSwitch failed: %v", err)
	}
	if ksResp.KillSwitchStatus != "DEACTIVATE" {
		t.Errorf("expected DEACTIVATE, got %s", ksResp.KillSwitchStatus)
	}

	// 3. PnL Exit
	pnlExit, err := client.GetPnLExit(ctx)
	if err != nil {
		t.Fatalf("GetPnLExit failed: %v", err)
	}
	expectedProfit := decimal.NewFromFloat(5000.50)
	if !pnlExit.Profit.Equal(expectedProfit) {
		t.Errorf("expected profit %v, got %v", expectedProfit, pnlExit.Profit)
	}

	// 4. SetPnLExit
	setPnLResp, err := client.SetPnLExit(ctx, models.PnlBasedExitRequest{
		ProfitValue:      decimal.NewFromFloat(10000),
		LossValue:        decimal.NewFromFloat(3000),
		EnableKillSwitch: true,
		ProductType:      []string{"INTRADAY"},
	})
	if err != nil {
		t.Fatalf("SetPnLExit failed: %v", err)
	}
	if setPnLResp.PnlExitStatus != "ACTIVE" {
		t.Errorf("expected ACTIVE, got %s", setPnLResp.PnlExitStatus)
	}

	// 5. StopPnLExit
	delResp, err := client.StopPnLExit(ctx)
	if err != nil {
		t.Fatalf("StopPnLExit failed: %v", err)
	}
	if delResp.Status != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", delResp.Status)
	}
}

func TestSecurity_ValidationNegativeCases(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Empty IP validation
	_, err := client.SetIP(ctx, models.UserIPRequest{IP: "", IPFlag: "PRIMARY"})
	if err == nil {
		t.Error("expected error for empty IP, got nil")
	}

	// Invalid KillSwitch status
	_, err = client.SetKillSwitch(ctx, "INVALID_STATUS")
	if err == nil {
		t.Error("expected error for invalid KillSwitch status, got nil")
	}
}

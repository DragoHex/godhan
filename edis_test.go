package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEDIS_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == URIEDISTpin:
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(`{"status":"SUCCESS","message":"TPIN generated and sent to registered mobile/email"}`))
		case r.Method == http.MethodPost && r.URL.Path == URIEDISForm:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"dhanClientId":"test-client","edisFormHtml":"<html><form id=\"edis\" action=\"https://edis.cdslindia.com\"></form></html>"}`))
		case r.Method == http.MethodPost && r.URL.Path == URIEDISBulkForm:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"dhanClientId":"test-client","edisFormHtml":"<html><form id=\"bulk-edis\"></form></html>"}`))
		case r.Method == http.MethodGet && r.URL.Path == URIEDISInquire+"/INE002A01018":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"clientId": "test-client",
				"isin": "INE002A01018",
				"totalQty": "100",
				"aprvdQty": "50",
				"status": "APPROVED",
				"remarks": "Partial approval"
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GenerateTPIN
	tpinResp, err := client.GenerateTPIN(ctx)
	if err != nil {
		t.Fatalf("GenerateTPIN failed: %v", err)
	}
	if tpinResp.Status != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", tpinResp.Status)
	}

	// 2. GenerateEDISForm
	formResp, err := client.GenerateEDISForm(ctx, models.EdisFormRequest{
		ISIN:     "INE002A01018",
		Qty:      10,
		Exchange: "NSE",
		Segment:  "EQ",
		Bulk:     false,
	})
	if err != nil {
		t.Fatalf("GenerateEDISForm failed: %v", err)
	}
	if formResp.EdisFormHTML == "" {
		t.Errorf("expected non-empty edisFormHtml")
	}

	// 3. GenerateBulkEDISForm
	bulkResp, err := client.GenerateBulkEDISForm(ctx, models.EdisBulkFormRequest{
		ISIN:     []string{"INE002A01018", "INE009A01021"},
		Exchange: "NSE",
		Segment:  "EQ",
	})
	if err != nil {
		t.Fatalf("GenerateBulkEDISForm failed: %v", err)
	}
	if bulkResp.EdisFormHTML == "" {
		t.Errorf("expected non-empty edisFormHtml")
	}

	// 4. InquireEDISQuantity
	qtyResp, err := client.InquireEDISQuantity(ctx, "INE002A01018")
	if err != nil {
		t.Fatalf("InquireEDISQuantity failed: %v", err)
	}
	if qtyResp.ISIN != "INE002A01018" || qtyResp.TotalQty != "100" || qtyResp.AprvdQty != "50" {
		t.Errorf("unexpected qty status response: %+v", qtyResp)
	}
}

func TestEDIS_Validation(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Inquire with empty ISIN
	_, err := client.InquireEDISQuantity(ctx, "")
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for empty ISIN, got %v", err)
	}

	// GenerateEDISForm with zero qty
	_, err = client.GenerateEDISForm(ctx, models.EdisFormRequest{
		ISIN: "INE002A01018",
		Qty:  0,
	})
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for zero qty, got %v", err)
	}
}

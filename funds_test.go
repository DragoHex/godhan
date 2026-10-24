package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestFunds_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == URIFundLimit:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"dhanClientId": "test-client",
				"availabelBalance": 105000.75,
				"sodLimit": 100000.00,
				"collateralAmount": 25000.00,
				"utilizedAmount": 15000.50,
				"withdrawableBalance": 85000.25
			}`))
		case r.Method == http.MethodPost && r.URL.Path == URIMarginCalculator:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"totalMargin": 4500.00,
				"spanMargin": 3000.00,
				"exposureMargin": 1500.00,
				"availableBalance": 105000.75,
				"insufficientBalance": 0.00
			}`))
		case r.Method == http.MethodPost && r.URL.Path == URIMarginCalculatorMulti:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"totalMargin": 7500.00,
				"spanMargin": 5000.00,
				"exposureMargin": 2500.00
			}`))
		case r.Method == http.MethodGet && r.URL.Path == URILedger:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"narration": "Funds Payin",
					"voucherdate": "2023-10-15",
					"credit": "50000.00",
					"debit": "0.00",
					"runbal": "105000.75"
				}
			]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GetFundLimits
	funds, err := client.GetFundLimits(ctx)
	if err != nil {
		t.Fatalf("GetFundLimits failed: %v", err)
	}
	expectedAvail := decimal.NewFromFloat(105000.75)
	if !funds.AvailableBalance.Equal(expectedAvail) {
		t.Errorf("expected balance %v, got %v", expectedAvail, funds.AvailableBalance)
	}

	// 2. CalculateMargin
	margin, err := client.CalculateMargin(ctx, models.KnowYourMarginReq{
		ExchangeSegment: "NSE_EQ",
		TransactionType: "BUY",
		Quantity:        100,
		ProductType:     "CNC",
		SecurityID:      "11536",
		Price:           decimal.NewFromFloat(3500.00),
	})
	if err != nil {
		t.Fatalf("CalculateMargin failed: %v", err)
	}
	expectedMargin := decimal.NewFromFloat(4500.00)
	if !margin.TotalMargin.Equal(expectedMargin) {
		t.Errorf("expected total margin %v, got %v", expectedMargin, margin.TotalMargin)
	}

	// 3. CalculateMultiMargin
	multiMargin, err := client.CalculateMultiMargin(ctx, models.MultiScripMarginCalcRequest{})
	if err != nil {
		t.Fatalf("CalculateMultiMargin failed: %v", err)
	}
	expectedMulti := decimal.NewFromFloat(7500.00)
	if !multiMargin.TotalMargin.Equal(expectedMulti) {
		t.Errorf("expected multi margin %v, got %v", expectedMulti, multiMargin.TotalMargin)
	}

	// 4. GetLedger
	ledger, err := client.GetLedger(ctx, "2023-10-01", "2023-10-15")
	if err != nil {
		t.Fatalf("GetLedger failed: %v", err)
	}
	if len(ledger) != 1 || ledger[0].Narration != "Funds Payin" {
		t.Errorf("unexpected ledger response: %+v", ledger)
	}
}

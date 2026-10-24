package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestPortfolio_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == URIHoldings:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"tradingSymbol": "TCS",
					"securityId": "11536",
					"isin": "INE467B01029",
					"totalQty": 25,
					"dpQty": 25,
					"t1Qty": 0,
					"availableQty": 25,
					"collateralQty": 0,
					"avgCostPrice": 3500.50
				}
			]`))
		case r.Method == http.MethodGet && r.URL.Path == URIPositions:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"tradingSymbol": "NIFTY-OCT-2023",
					"securityId": "52175",
					"positionType": "LONG",
					"exchangeSegment": "NSE_FNO",
					"productType": "INTRADAY",
					"buyAvg": 19500.25,
					"buyQty": 50,
					"sellAvg": 0.0,
					"sellQty": 0,
					"netQty": 50,
					"realizedProfit": 0.0,
					"unrealizedProfit": 1250.75
				}
			]`))
		case r.Method == http.MethodPost && r.URL.Path == URIPositionsConvert:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"SUCCESS","message":"models.Position converted successfully"}`))
		case r.Method == http.MethodDelete && r.URL.Path == URIPositions:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"SUCCESS","message":"All positions exited"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GetHoldings
	holdings, err := client.GetHoldings(ctx)
	if err != nil {
		t.Fatalf("GetHoldings failed: %v", err)
	}
	if len(holdings) != 1 || holdings[0].TradingSymbol != "TCS" || holdings[0].TotalQty != 25 {
		t.Errorf("unexpected holdings: %+v", holdings)
	}
	expectedHoldingPrice := decimal.NewFromFloat(3500.50)
	if !holdings[0].AvgCostPrice.Equal(expectedHoldingPrice) {
		t.Errorf("expected avg price %v, got %v", expectedHoldingPrice, holdings[0].AvgCostPrice)
	}

	// 2. GetPositions
	positions, err := client.GetPositions(ctx)
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}
	if len(positions) != 1 || positions[0].NetQty != 50 {
		t.Errorf("unexpected positions: %+v", positions)
	}

	// 3. ConvertPosition
	err = client.ConvertPosition(ctx, models.PositionConversionRequest{
		FromProductType: "INTRADAY",
		ToProductType:   "CNC",
		ExchangeSegment: "NSE_EQ",
		PositionType:    "LONG",
		SecurityID:      "11536",
		ConvertQty:      10,
	})
	if err != nil {
		t.Fatalf("ConvertPosition failed: %v", err)
	}

	// 4. ExitAllPositions
	exitResp, err := client.ExitAllPositions(ctx)
	if err != nil {
		t.Fatalf("ExitAllPositions failed: %v", err)
	}
	if exitResp.Status != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", exitResp.Status)
	}
}

func TestPortfolio_ValidationNegative(t *testing.T) {
	client := NewClient("test-client", "test-token")
	ctx := context.Background()

	// Convert zero quantity
	err := client.ConvertPosition(ctx, models.PositionConversionRequest{
		FromProductType: "INTRADAY",
		ToProductType:   "CNC",
		ConvertQty:      0,
	})
	if err == nil || !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

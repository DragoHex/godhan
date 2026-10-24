package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCharts_Lifecycle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == URIChartsHistorical:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"open": [100.0, 102.5],
				"high": [105.0, 107.0],
				"low": [98.0, 101.0],
				"close": [102.0, 106.5],
				"volume": [10000, 15000],
				"timestamp": [1697356800, 1697443200],
				"open_interest": [500, 600]
			}`))
		case r.Method == http.MethodPost && r.URL.Path == URIChartsIntraday:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"open": [100.0],
				"high": [101.0],
				"low": [99.5],
				"close": [100.5],
				"volume": [2000],
				"timestamp": [1697357000],
				"open_interest": [100]
			}`))
		case r.Method == http.MethodPost && r.URL.Path == URIChartsRollingOption:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"open": [25.0],
				"high": [30.0],
				"low": [20.0],
				"close": [28.0],
				"volume": [5000],
				"timestamp": [1697357000]
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. GetHistoricalCharts
	histCharts, err := client.GetHistoricalCharts(ctx, models.HistoricalChartsRequest{
		SecurityID:      "1333",
		ExchangeSegment: "NSE_EQ",
		Instrument:      "EQUITY",
		FromDate:        "2023-10-01",
		ToDate:          "2023-10-15",
	})
	if err != nil {
		t.Fatalf("GetHistoricalCharts failed: %v", err)
	}
	if len(histCharts.Open) != 2 {
		t.Fatalf("expected 2 candles, got %d", len(histCharts.Open))
	}

	// 2. ToCandles conversion helper (Kite parity)
	candles := histCharts.ToCandles()
	if len(candles) != 2 {
		t.Fatalf("expected 2 candles from ToCandles, got %d", len(candles))
	}
	expectedOpen := decimal.NewFromFloat(100.0)
	if !candles[0].Open.Equal(expectedOpen) {
		t.Errorf("expected candle[0].Open %v, got %v", expectedOpen, candles[0].Open)
	}

	// 3. GetIntradayCharts
	intraCharts, err := client.GetIntradayCharts(ctx, models.IntradayChartsRequest{
		SecurityID:      "1333",
		ExchangeSegment: "NSE_EQ",
		Instrument:      "EQUITY",
		Interval:        "1",
		FromDate:        "2023-10-15 09:15:00",
		ToDate:          "2023-10-15 15:30:00",
	})
	if err != nil {
		t.Fatalf("GetIntradayCharts failed: %v", err)
	}
	if len(intraCharts.Open) != 1 {
		t.Fatalf("expected 1 candle, got %d", len(intraCharts.Open))
	}

	// 4. GetRollingOptionCharts
	rollingCharts, err := client.GetRollingOptionCharts(ctx, models.OptionChartRequest{})
	if err != nil {
		t.Fatalf("GetRollingOptionCharts failed: %v", err)
	}
	if len(rollingCharts.Open) != 1 {
		t.Fatalf("expected 1 rolling candle, got %d", len(rollingCharts.Open))
	}
}

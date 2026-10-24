package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestData_AnalyticsAndMaster(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == URIDataCompanyInfo:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"companyName":"Tata Consultancy Services Ltd"}`))
		case r.Method == http.MethodPost && r.URL.Path == URIDataMarketMovers:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"gainers":[{"tradingSymbol":"RELIANCE"}],"losers":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == URIDataNewsHeadline:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"headlines":[{"title":"Q2 Earnings Announcement"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == URIDataTechnical:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"sma50":3450.0,"ema20":3480.0}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient("test-client", "test-token", WithBaseURI(ts.URL))
	ctx := context.Background()

	// 1. CompanyInfo
	comp, err := client.GetCompanyInfo(ctx, models.FundamentalMetricsRequest{SecurityID: "11536"})
	if err != nil {
		t.Fatalf("GetCompanyInfo failed: %v", err)
	}
	if comp.CompanyName != "Tata Consultancy Services Ltd" {
		t.Errorf("unexpected company info: %+v", comp)
	}

	// 2. MarketMovers
	movers, err := client.GetMarketMovers(ctx, models.TopInstrumentsRequest{})
	if err != nil {
		t.Fatalf("GetMarketMovers failed: %v", err)
	}
	if len(movers.Gainers) != 1 {
		t.Errorf("unexpected movers: %+v", movers)
	}

	// 3. NewsHeadlines
	news, err := client.GetNewsHeadlines(ctx, models.NewsHeadlineRequest{})
	if err != nil {
		t.Fatalf("GetNewsHeadlines failed: %v", err)
	}
	if len(news.Headlines) != 1 {
		t.Errorf("unexpected headlines: %+v", news)
	}

	// 4. TechnicalData
	tech, err := client.GetTechnicalData(ctx, models.TechnicalMetricsRequest{})
	if err != nil {
		t.Fatalf("GetTechnicalData failed: %v", err)
	}
	if tech == nil {
		t.Fatal("expected non-nil technical data")
	}

	// 5. Streaming CSV models.Instrument Master (KiteSDK parity)
	sampleCSV := `SEM_EXM_EXCH_ID,SEM_SEGMENT,SEM_SMST_SECURITY_ID,SEM_INSTRUMENT_NAME,SEM_EXPIRY_DATE,SEM_TRADING_SYMBOL,SEM_LOT_UNITS,SEM_CUSTOM_SYMBOL,SEM_STRIKE_PRICE
NSE,E,11536,EQUITY,,TCS-EQ,1,TCS,0.00
NSE,D,52175,FUTIDX,2023-10-26,NIFTY-26OCT23-FUT,50,NIFTY,0.00
`
	instruments, err := ParseInstruments(strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("ParseInstruments failed: %v", err)
	}
	if len(instruments) != 2 {
		t.Fatalf("expected 2 instruments, got %d", len(instruments))
	}
	if instruments[0].TradingSymbol != "TCS-EQ" || instruments[0].LotSize != 1 {
		t.Errorf("unexpected instrument[0]: %+v", instruments[0])
	}
	if instruments[1].SecurityID != "52175" || instruments[1].LotSize != 50 {
		t.Errorf("unexpected instrument[1]: %+v", instruments[1])
	}
}

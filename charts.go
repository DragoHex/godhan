package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"

)

// models.HistoricalChartsRequest specifies parameters to fetch daily OHLC candle arrays

// models.IntradayChartsRequest specifies parameters to fetch minute-level OHLC candle arrays

// models.OptionChartRequest specifies parameters to fetch expired options historical data

// models.ChartsResponse represents parallel arrays of OHLC historical market data from DhanHQ

// models.OptionChartResponse returns expired option chart data

// models.Candle represents an ergonomic, strongly-typed individual OHLC bar (KiteSDK parity)



// GetHistoricalCharts queries daily historical OHLC candle data
func (c *Client) GetHistoricalCharts(ctx context.Context, req models.HistoricalChartsRequest) (*models.ChartsResponse, error) {
	var resp models.ChartsResponse
	if err := c.doRequest(ctx, http.MethodPost, URIChartsHistorical, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetIntradayCharts queries minute-level historical OHLC candle data
func (c *Client) GetIntradayCharts(ctx context.Context, req models.IntradayChartsRequest) (*models.ChartsResponse, error) {
	var resp models.ChartsResponse
	if err := c.doRequest(ctx, http.MethodPost, URIChartsIntraday, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRollingOptionCharts queries historical candles for expired options
func (c *Client) GetRollingOptionCharts(ctx context.Context, req models.OptionChartRequest) (*models.OptionChartResponse, error) {
	var resp models.OptionChartResponse
	if err := c.doRequest(ctx, http.MethodPost, URIChartsRollingOption, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

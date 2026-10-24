package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
)

// models.FundamentalMetricsRequest specifies scrip to fetch company information

// models.FundamentalMetricsResponse contains company fundamental details

// models.TopInstrumentsRequest specifies criteria to fetch top gainers and losers

// models.TopInstrumentItem details one market mover

// models.TopInstrumentsResponse contains top gainers and losers

// models.NewsHeadlineRequest queries news headlines for symbol

// models.NewsItem details a news story

// models.NewsHeadlineResponse returns news headlines

// models.TechnicalMetricsRequest queries technical indicators

// models.TechnicalMetricsResponse returns technical indicators

// GetCompanyInfo fetches corporate and fundamental details for an equity scrip
func (c *Client) GetCompanyInfo(ctx context.Context, req models.FundamentalMetricsRequest) (*models.FundamentalMetricsResponse, error) {
	var resp models.FundamentalMetricsResponse
	if err := c.doRequest(ctx, http.MethodPost, URIDataCompanyInfo, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMarketMovers fetches top gaining and losing instruments
func (c *Client) GetMarketMovers(ctx context.Context, req models.TopInstrumentsRequest) (*models.TopInstrumentsResponse, error) {
	var resp models.TopInstrumentsResponse
	if err := c.doRequest(ctx, http.MethodPost, URIDataMarketMovers, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNewsHeadlines fetches financial news headlines
func (c *Client) GetNewsHeadlines(ctx context.Context, req models.NewsHeadlineRequest) (*models.NewsHeadlineResponse, error) {
	var resp models.NewsHeadlineResponse
	if err := c.doRequest(ctx, http.MethodPost, URIDataNewsHeadline, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTechnicalData fetches calculated technical indicators
func (c *Client) GetTechnicalData(ctx context.Context, req models.TechnicalMetricsRequest) (*models.TechnicalMetricsResponse, error) {
	var resp models.TechnicalMetricsResponse
	if err := c.doRequest(ctx, http.MethodPost, URIDataTechnical, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

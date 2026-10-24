package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"

)

// models.FundLimitResponse details capital, cash balances, collateral, and margin utilization

// models.KnowYourMarginReq specifies order parameters to query pre-trade margin requirement

// models.KnowYourMarginResponse returns breakdown of required margin for a single order

// models.MultiScripMarginCalcRequest specifies multi-leg basket of orders for portfolio margin calculation

// models.MultiScripMarginCalcResponse returns composite margin requirements for a multi-leg basket

// GetFundLimits fetches available trading balances and collateral limits
func (c *Client) GetFundLimits(ctx context.Context) (*models.FundLimitResponse, error) {
	var resp models.FundLimitResponse
	if err := c.doRequest(ctx, http.MethodGet, URIFundLimit, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CalculateMargin computes pre-trade margin requirement for a single order
func (c *Client) CalculateMargin(ctx context.Context, req models.KnowYourMarginReq) (*models.KnowYourMarginResponse, error) {
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp models.KnowYourMarginResponse
	if err := c.doRequest(ctx, http.MethodPost, URIMarginCalculator, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CalculateMultiMargin computes basket margin and hedging benefits across multiple legs
func (c *Client) CalculateMultiMargin(ctx context.Context, req models.MultiScripMarginCalcRequest) (*models.MultiScripMarginCalcResponse, error) {
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp models.MultiScripMarginCalcResponse
	if err := c.doRequest(ctx, http.MethodPost, URIMarginCalculatorMulti, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

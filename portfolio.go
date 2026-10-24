package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"

)

// models.Holding represents a Demat equity custody holding

// models.Position represents an open intraday or carry-forward derivative position

// models.PositionConversionRequest payload to switch product type (e.g. INTRADAY <-> CNC / MARGIN)

// GetHoldings retrieves long-term Demat holdings
func (c *Client) GetHoldings(ctx context.Context) ([]models.Holding, error) {
	var resp []models.Holding
	if err := c.doRequest(ctx, http.MethodGet, URIHoldings, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPositions retrieves open and closed positions for the trading day
func (c *Client) GetPositions(ctx context.Context) ([]models.Position, error) {
	var resp []models.Position
	if err := c.doRequest(ctx, http.MethodGet, URIPositions, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ConvertPosition switches position product type
func (c *Client) ConvertPosition(ctx context.Context, req models.PositionConversionRequest) error {
	if req.ConvertQty <= 0 {
		return fmt.Errorf("%w: convertQty must be greater than zero", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	return c.doRequest(ctx, http.MethodPost, URIPositionsConvert, nil, req, nil)
}

// ExitAllPositions squares off all open intraday and derivative positions
func (c *Client) ExitAllPositions(ctx context.Context) (*models.UserIPResponse, error) {
	var resp models.UserIPResponse
	if err := c.doRequest(ctx, http.MethodDelete, URIPositions, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

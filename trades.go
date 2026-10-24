package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"
	"strings"

)

// models.TradeResponse details an executed trade fill

// GetTrades retrieves all executed trades for the current trading day
func (c *Client) GetTrades(ctx context.Context) ([]models.TradeResponse, error) {
	var resp []models.TradeResponse
	if err := c.doRequest(ctx, http.MethodGet, URITrades, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTradesByOrderID retrieves all trade fills associated with a specific order
func (c *Client) GetTradesByOrderID(ctx context.Context, orderID string) ([]models.TradeResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URITrades, cleanID)
	var resp []models.TradeResponse
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTradeHistory retrieves historical trade executions within a specified date window
func (c *Client) GetTradeHistory(ctx context.Context, fromDate, toDate string, page int) ([]models.TradeResponse, error) {
	cleanFrom := strings.TrimSpace(fromDate)
	cleanTo := strings.TrimSpace(toDate)
	if cleanFrom == "" || cleanTo == "" {
		return nil, fmt.Errorf("%w: fromDate and toDate must be provided", ErrInvalidInput)
	}
	if page < 1 {
		return nil, fmt.Errorf("%w: page number must be >= 1", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s/%s/%d", URITrades, cleanFrom, cleanTo, page)
	var resp []models.TradeResponse
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetOrderTrades retrieves all trades executed for a given order (parity alias for GetTradesByOrderID)
func (c *Client) GetOrderTrades(ctx context.Context, orderID string) ([]models.TradeResponse, error) {
	return c.GetTradesByOrderID(ctx, orderID)
}

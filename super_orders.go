package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"
	"strings"

)

// models.SuperOrderRequest specifies parameters for bracket orders (entry, stop-loss, and target legs)

// models.SuperModifyRequest specifies parameters to update pending super order legs

// models.SuperOrderResponse represents super order state from super order book

// PlaceSuperOrder places a bracket order with automated target and stop loss legs
func (c *Client) PlaceSuperOrder(ctx context.Context, req models.SuperOrderRequest) (*models.OrderStatusResponse, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.ClientID()
	}
	var resp models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPost, URISuperOrders, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifySuperOrder modifies target, stop loss, or entry price of a super order
func (c *Client) ModifySuperOrder(ctx context.Context, orderID string, req models.SuperModifyRequest) (*models.OrderStatusResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.ClientID()
	}
	req.OrderID = cleanID

	endpoint := fmt.Sprintf("%s/%s", URISuperOrders, cleanID)
	var resp models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPut, endpoint, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelSuperOrder cancels an individual super order leg or entire bracket order
func (c *Client) CancelSuperOrder(ctx context.Context, orderID string, orderLeg string) (*models.OrderStatusResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	cleanLeg := strings.TrimSpace(orderLeg)
	if cleanID == "" || cleanLeg == "" {
		return nil, fmt.Errorf("%w: orderId and orderLeg must both be provided", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s/%s", URISuperOrders, cleanID, cleanLeg)
	var resp models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSuperOrders retrieves all super/bracket orders placed today
func (c *Client) GetSuperOrders(ctx context.Context) ([]models.SuperOrderResponse, error) {
	var resp []models.SuperOrderResponse
	if err := c.doRequest(ctx, http.MethodGet, URISuperOrders, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

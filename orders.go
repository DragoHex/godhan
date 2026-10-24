package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"
	"strings"

)

// models.OrderRequest specifies parameters to place an order

// models.OrderModifyRequest specifies fields that can be updated on a pending order

// models.OrderStatusResponse contains order placement/modification/cancellation status

// models.OrderResponse represents full order audit details from order book

// PlaceOrder places a regular, AMO, or intraday order on DhanHQ
func (c *Client) PlaceOrder(ctx context.Context, req models.OrderRequest) (*models.OrderStatusResponse, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPost, URIOrders, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyOrder modifies price, trigger price, or quantity of a pending order
func (c *Client) ModifyOrder(ctx context.Context, orderID string, req models.OrderModifyRequest) (*models.OrderStatusResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	req.OrderID = cleanID

	endpoint := fmt.Sprintf("%s/%s", URIOrders, cleanID)
	var resp models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPut, endpoint, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelOrder cancels a pending order
func (c *Client) CancelOrder(ctx context.Context, orderID string) (*models.OrderStatusResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIOrders, cleanID)
	var resp models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrders retrieves all orders placed during the current trading day
func (c *Client) GetOrders(ctx context.Context) ([]models.OrderResponse, error) {
	var resp []models.OrderResponse
	if err := c.doRequest(ctx, http.MethodGet, URIOrders, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetOrderByID retrieves specific order details by internal Dhan order ID
func (c *Client) GetOrderByID(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIOrders, cleanID)
	var resp models.OrderResponse
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrderByCorrelationID retrieves order by client external correlation ID
func (c *Client) GetOrderByCorrelationID(ctx context.Context, correlationID string) (*models.OrderResponse, error) {
	cleanID := strings.TrimSpace(correlationID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: correlationId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIOrdersExternal, cleanID)
	var resp models.OrderResponse
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SliceOrder auto-splits large quantity orders across exchange freeze limits
func (c *Client) SliceOrder(ctx context.Context, req models.OrderRequest) ([]models.OrderStatusResponse, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp []models.OrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPost, URIOrdersSlicing, nil, req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ExitOrder cancels an open or pending order (parity alias for CancelOrder)
func (c *Client) ExitOrder(ctx context.Context, orderID string) (*models.OrderStatusResponse, error) {
	return c.CancelOrder(ctx, orderID)
}

// GetOrderHistory retrieves audit history for a specific order (parity alias for GetOrderByID)
func (c *Client) GetOrderHistory(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	return c.GetOrderByID(ctx, orderID)
}

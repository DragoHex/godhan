package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"
	"strings"

)

// models.GTTOrderModel represents a Good-Til-Triggered / Forever order payload

// models.GttModifyRequest updates trigger or limit prices for a pending forever order

// models.GttOrderStatusResponse confirms forever order submission or cancellation

// models.ForeverOrderResponse details active forever order

// models.AlertCondition defines criteria to trigger conditional orders

// models.AlertOrderRequest payload for creating conditional alerts and execution triggers

// models.AlertModifyRequest payload for modifying conditional alerts

// models.AlertOrderResponse confirms alert creation

// models.GetAlertResponse returns detailed conditional alert record

// models.MultiOrderRequest groups several orders to be placed simultaneously

// models.MultiOrderResponse returns execution status of batch order

// PlaceForeverOrder creates a GTT (Good-Til-Triggered) conditional order
func (c *Client) PlaceForeverOrder(ctx context.Context, req models.GTTOrderModel) (*models.GttOrderStatusResponse, error) {
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp models.GttOrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPost, URIForeverOrders, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyForeverOrder updates price or trigger parameters of an open GTT order
func (c *Client) ModifyForeverOrder(ctx context.Context, orderID string, req models.GttModifyRequest) (*models.GttOrderStatusResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	req.OrderID = cleanID

	endpoint := fmt.Sprintf("%s/%s", URIForeverOrders, cleanID)
	var resp models.GttOrderStatusResponse
	if err := c.doRequest(ctx, http.MethodPut, endpoint, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelForeverOrder cancels an open GTT order
func (c *Client) CancelForeverOrder(ctx context.Context, orderID string) (*models.GttOrderStatusResponse, error) {
	cleanID := strings.TrimSpace(orderID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: orderId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIForeverOrders, cleanID)
	var resp models.GttOrderStatusResponse
	if err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetForeverOrders retrieves all active and historical GTT orders
func (c *Client) GetForeverOrders(ctx context.Context) ([]models.ForeverOrderResponse, error) {
	var resp []models.ForeverOrderResponse
	if err := c.doRequest(ctx, http.MethodGet, URIForeverOrders, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// PlaceAlertOrder creates conditional trigger alert
func (c *Client) PlaceAlertOrder(ctx context.Context, req models.AlertOrderRequest) (*models.AlertOrderResponse, error) {
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp models.AlertOrderResponse
	if err := c.doRequest(ctx, http.MethodPost, URIAlertOrders, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAlertOrders lists all configured conditional alerts
func (c *Client) GetAlertOrders(ctx context.Context) ([]models.GetAlertResponse, error) {
	var resp []models.GetAlertResponse
	if err := c.doRequest(ctx, http.MethodGet, URIAlertOrders, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAlertOrderByID fetches conditional alert details by alert ID
func (c *Client) GetAlertOrderByID(ctx context.Context, alertID string) (*models.GetAlertResponse, error) {
	cleanID := strings.TrimSpace(alertID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: alertId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIAlertOrders, cleanID)
	var resp models.GetAlertResponse
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyAlertOrder modifies conditional alert thresholds
func (c *Client) ModifyAlertOrder(ctx context.Context, alertID string, req models.AlertModifyRequest) (*models.AlertOrderResponse, error) {
	cleanID := strings.TrimSpace(alertID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: alertId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIAlertOrders, cleanID)
	var resp models.AlertOrderResponse
	if err := c.doRequest(ctx, http.MethodPut, endpoint, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelAlertOrder removes a conditional alert
func (c *Client) CancelAlertOrder(ctx context.Context, alertID string) (*models.AlertOrderResponse, error) {
	cleanID := strings.TrimSpace(alertID)
	if cleanID == "" {
		return nil, fmt.Errorf("%w: alertId cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIAlertOrders, cleanID)
	var resp models.AlertOrderResponse
	if err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PlaceMultiOrder dispatches a batch of orders simultaneously
func (c *Client) PlaceMultiOrder(ctx context.Context, req models.MultiOrderRequest) (*models.MultiOrderResponse, error) {
	if req.DhanClientID == "" {
		req.DhanClientID = c.clientID
	}
	var resp models.MultiOrderResponse
	if err := c.doRequest(ctx, http.MethodPost, URIMultiAlertOrders, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

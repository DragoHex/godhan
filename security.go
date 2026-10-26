package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

)

// models.UserIPRequest defines the payload for setting or modifying static IP

// models.UserIPResponse represents standard status response for IP and exit endpoints

// models.GetIPDetailsResponse details currently registered static IPs and matching status

// models.KillSwitchResponse represents account kill switch activation status

// models.PnlBasedExitRequest payload to configure daily P&L square-off limits

// models.PnlExitResponse confirms configuration of P&L exit

// models.ExitPnlResponse returns active P&L exit parameters

// GetIP fetches registered static IP configuration
func (c *Client) GetIP(ctx context.Context) (*models.GetIPDetailsResponse, error) {
	var resp models.GetIPDetailsResponse
	if err := c.doRequest(ctx, http.MethodGet, URIGetIP, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetIP registers a new static IP (PRIMARY or SECONDARY)
func (c *Client) SetIP(ctx context.Context, req models.UserIPRequest) (*models.UserIPResponse, error) {
	if strings.TrimSpace(req.IP) == "" {
		return nil, fmt.Errorf("%w: IP cannot be empty", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.ClientID()
	}
	var resp models.UserIPResponse
	if err := c.doRequest(ctx, http.MethodPost, URISetIP, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyIP updates an existing static IP
func (c *Client) ModifyIP(ctx context.Context, req models.UserIPRequest) (*models.UserIPResponse, error) {
	if strings.TrimSpace(req.IP) == "" {
		return nil, fmt.Errorf("%w: IP cannot be empty", ErrInvalidInput)
	}
	if req.DhanClientID == "" {
		req.DhanClientID = c.ClientID()
	}
	var resp models.UserIPResponse
	if err := c.doRequest(ctx, http.MethodPut, URIModifyIP, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetKillSwitchStatus checks whether the account trading kill switch is active
func (c *Client) GetKillSwitchStatus(ctx context.Context) (*models.KillSwitchResponse, error) {
	var resp models.KillSwitchResponse
	if err := c.doRequest(ctx, http.MethodGet, URIKillSwitch, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetKillSwitch activates or deactivates the trading kill switch ("ACTIVATE" or "DEACTIVATE")
func (c *Client) SetKillSwitch(ctx context.Context, status string) (*models.KillSwitchResponse, error) {
	upperStatus := strings.ToUpper(strings.TrimSpace(status))
	if upperStatus != "ACTIVATE" && upperStatus != "DEACTIVATE" {
		return nil, fmt.Errorf("%w: killSwitchStatus must be ACTIVATE or DEACTIVATE", ErrInvalidInput)
	}

	query := url.Values{}
	query.Set("killSwitchStatus", upperStatus)

	var resp models.KillSwitchResponse
	if err := c.doRequest(ctx, http.MethodPost, URIKillSwitch, query, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPnLExit fetches current P&L based exit configurations
func (c *Client) GetPnLExit(ctx context.Context) (*models.ExitPnlResponse, error) {
	var resp models.ExitPnlResponse
	if err := c.doRequest(ctx, http.MethodGet, URIPnLExit, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetPnLExit configures automated P&L thresholds to exit open positions
func (c *Client) SetPnLExit(ctx context.Context, req models.PnlBasedExitRequest) (*models.PnlExitResponse, error) {
	if req.DhanClientID == "" {
		req.DhanClientID = c.ClientID()
	}
	var resp models.PnlExitResponse
	if err := c.doRequest(ctx, http.MethodPost, URIPnLExit, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StopPnLExit cancels active P&L monitoring and square-off rule
func (c *Client) StopPnLExit(ctx context.Context) (*models.UserIPResponse, error) {
	var resp models.UserIPResponse
	if err := c.doRequest(ctx, http.MethodDelete, URIPnLExit, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

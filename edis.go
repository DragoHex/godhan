package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// models.TPINResponse captures T-PIN generation confirmation

// models.EdisFormRequest defines payload for single-stock electronic delivery authorization

// models.EdisBulkFormRequest defines payload for multi-stock electronic delivery authorization

// models.EdisFormResponse returns generated CDSL redirect HTML form

// models.EdisQtyStatusResponse provides authorized vs allocated share quantity for ISIN

// GenerateTPIN requests CDSL to send a new T-PIN to the client's registered mobile/email
func (c *Client) GenerateTPIN(ctx context.Context) (*models.TPINResponse, error) {
	var resp models.TPINResponse
	if err := c.doRequest(ctx, http.MethodGet, URIEDISTpin, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GenerateEDISForm generates HTML form payload for CDSL T-PIN authorization of a single scrip
func (c *Client) GenerateEDISForm(ctx context.Context, req models.EdisFormRequest) (*models.EdisFormResponse, error) {
	if strings.TrimSpace(req.ISIN) == "" {
		return nil, fmt.Errorf("%w: ISIN cannot be empty", ErrInvalidInput)
	}
	if req.Qty <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidInput)
	}
	var resp models.EdisFormResponse
	if err := c.doRequest(ctx, http.MethodPost, URIEDISForm, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GenerateBulkEDISForm generates HTML form payload for CDSL T-PIN authorization of multiple scrips
func (c *Client) GenerateBulkEDISForm(ctx context.Context, req models.EdisBulkFormRequest) (*models.EdisFormResponse, error) {
	if len(req.ISIN) == 0 {
		return nil, fmt.Errorf("%w: ISIN list cannot be empty", ErrInvalidInput)
	}
	var resp models.EdisFormResponse
	if err := c.doRequest(ctx, http.MethodPost, URIEDISBulkForm, nil, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// InquireEDISQuantity checks current EDIS authorization status and approved quantity for an ISIN
func (c *Client) InquireEDISQuantity(ctx context.Context, isin string) (*models.EdisQtyStatusResponse, error) {
	cleanISIN := strings.TrimSpace(isin)
	if cleanISIN == "" {
		return nil, fmt.Errorf("%w: ISIN cannot be empty", ErrInvalidInput)
	}
	endpoint := fmt.Sprintf("%s/%s", URIEDISInquire, cleanISIN)
	var resp models.EdisQtyStatusResponse
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

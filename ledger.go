package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
	"net/url"
)

// models.BoLedgerEntry represents an individual accounting ledger transaction record

// GetLedger retrieves statement of ledger financial transactions
func (c *Client) GetLedger(ctx context.Context, fromDate, toDate string) ([]models.BoLedgerEntry, error) {
	query := url.Values{}
	if fromDate != "" {
		query.Set("from-date", fromDate)
	}
	if toDate != "" {
		query.Set("to-date", toDate)
	}

	var resp []models.BoLedgerEntry
	if err := c.doRequest(ctx, http.MethodGet, URILedger, query, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

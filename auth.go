package godhan

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/DragoHex/godhan/models"
)

// RenewToken refreshes an active access token for an additional 24 hours.
// If tokenOverride is passed, that token will be renewed; otherwise, c.AccessToken() is used.
// On success, the client's internal accessToken is automatically updated thread-safely.
func (c *Client) RenewToken(ctx context.Context, tokenOverride ...string) (*models.RenewTokenResponse, error) {
	origToken := c.AccessToken()
	token := origToken
	if len(tokenOverride) > 0 {
		token = strings.TrimSpace(tokenOverride[0])
	}

	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("%w: cannot renew empty access token", ErrInvalidInput)
	}

	// Temporarily apply token for the request
	c.SetAccessToken(token)

	var resp models.RenewTokenResponse
	err := c.doRequest(ctx, http.MethodGet, URIRenewToken, nil, nil, &resp)
	if err != nil {
		if len(tokenOverride) > 0 {
			c.SetAccessToken(origToken)
		}
		return nil, err
	}

	if resp.AccessToken != "" {
		c.SetAccessToken(resp.AccessToken)
	} else if len(tokenOverride) > 0 {
		c.SetAccessToken(origToken)
	}

	if resp.DhanClientID != "" && c.ClientID() == "" {
		c.SetClientID(resp.DhanClientID)
	}

	return &resp, nil
}

// RenewAccessToken is a parity alias for RenewToken.
func (c *Client) RenewAccessToken(ctx context.Context, tokenOverride ...string) (*models.RenewTokenResponse, error) {
	return c.RenewToken(ctx, tokenOverride...)
}

package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"net/http"
)

// models.UserProfile represents user profile and IP security details from DhanHQ

// GetUserProfile retrieves current user profile and whitelisted IP details
func (c *Client) GetUserProfile(ctx context.Context) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := c.doRequest(ctx, http.MethodGet, URIGetIP, nil, nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

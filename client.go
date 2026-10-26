package godhan

import "sync"

// Client is the primary interface for communicating with DhanHQ REST API
type Client struct {
	mu          sync.RWMutex
	clientID    string
	accessToken string
	baseURI     string
	httpClient  HTTPClient
	configErr   error
}

// Option represents a functional configuration option for Client
type Option func(*Client)

// WithBaseURI overrides the default DhanHQ REST API base URL
func WithBaseURI(uri string) Option {
	return func(c *Client) {
		if uri != "" {
			c.baseURI = uri
		}
	}
}

// WithHTTPClient provides a custom HTTP client (useful for fakes, proxy rotators, rate limiters)
func WithHTTPClient(client HTTPClient) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithAccessToken sets or overrides the access token
func WithAccessToken(token string) Option {
	return func(c *Client) {
		c.accessToken = token
	}
}

// WithClientID sets or overrides the client ID
func WithClientID(clientID string) Option {
	return func(c *Client) {
		c.clientID = clientID
	}
}

// NewClient creates a new configured DhanHQ API Client
func NewClient(clientID, accessToken string, opts ...Option) *Client {
	c := &Client{
		clientID:    clientID,
		accessToken: accessToken,
		baseURI:     DefaultBaseURI,
		httpClient:  defaultHTTPClient(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetClientID updates the configured Dhan Client ID thread-safely
func (c *Client) SetClientID(clientID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clientID = clientID
}

// ClientID returns the configured Dhan Client ID thread-safely
func (c *Client) ClientID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.clientID
}

// SetAccessToken updates the configured access token thread-safely
func (c *Client) SetAccessToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = token
}

// AccessToken returns the configured Access Token thread-safely
func (c *Client) AccessToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.accessToken
}

// New creates a new Client configured via functional options (parity alias)
func New(opts ...Option) *Client {
	c := &Client{
		baseURI:    DefaultBaseURI,
		httpClient: defaultHTTPClient(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// GetLoginURL generates the DhanHQ consent login URL for partner authentication
func (c *Client) GetLoginURL(partnerID string) string {
	clientID := c.ClientID()
	if partnerID != "" {
		return "https://auth.dhan.co/login/consent?partner_id=" + partnerID + "&client_id=" + clientID
	}
	return "https://auth.dhan.co/login/consent?client_id=" + clientID
}

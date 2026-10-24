package godhan

// Client is the primary interface for communicating with DhanHQ REST API
type Client struct {
	clientID    string
	accessToken string
	baseURI     string
	httpClient  HTTPClient
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

// ClientID returns the configured Dhan Client ID
func (c *Client) ClientID() string {
	return c.clientID
}

// AccessToken returns the configured Access Token
func (c *Client) AccessToken() string {
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
	if partnerID != "" {
		return "https://auth.dhan.co/login/consent?partner_id=" + partnerID + "&client_id=" + c.clientID
	}
	return "https://auth.dhan.co/login/consent?client_id=" + c.clientID
}

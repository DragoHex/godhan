package godhan

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ProxyConfig holds forward proxy credentials and endpoint for static IP routing.
type ProxyConfig struct {
	Scheme       string // "http" or "https" (default "https")
	Host         string // e.g. "dc46-mum-01.algoip.in" or "103.21.54.1"
	Port         int    // e.g. 443 (default 443 for https, 80 for http)
	ClientID     string // proxy username / client ID
	ClientSecret string // proxy password / client secret
}

// URL constructs a validated *url.URL from the ProxyConfig.
func (p ProxyConfig) URL() (*url.URL, error) {
	rawHost := strings.TrimSpace(p.Host)
	if rawHost == "" {
		return nil, fmt.Errorf("%w: proxy host cannot be empty", ErrInvalidInput)
	}

	// Strip protocol prefix if accidentally included in Host
	scheme := strings.ToLower(strings.TrimSpace(p.Scheme))
	if strings.HasPrefix(rawHost, "http://") {
		rawHost = strings.TrimPrefix(rawHost, "http://")
		if scheme == "" {
			scheme = "http"
		}
	} else if strings.HasPrefix(rawHost, "https://") {
		rawHost = strings.TrimPrefix(rawHost, "https://")
		if scheme == "" {
			scheme = "https"
		}
	}

	if scheme == "" {
		scheme = "https"
	}
	if scheme != "http" && scheme != "https" {
		return nil, fmt.Errorf("%w: unsupported proxy scheme %q (must be http or https)", ErrInvalidInput, p.Scheme)
	}

	host := rawHost
	port := p.Port

	// Check if port was already embedded in Host (e.g. host:port)
	if h, portStr, err := net.SplitHostPort(rawHost); err == nil {
		host = h
		if parsedPort, err := strconv.Atoi(portStr); err == nil && parsedPort > 0 {
			port = parsedPort
		}
	}

	if port <= 0 {
		if scheme == "http" {
			port = 80
		} else {
			port = 443
		}
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
	}

	if p.ClientID != "" {
		if p.ClientSecret != "" {
			u.User = url.UserPassword(p.ClientID, p.ClientSecret)
		} else {
			u.User = url.User(p.ClientID)
		}
	}

	return u, nil
}

// WithProxy configures an authenticated HTTP/HTTPS forward proxy on the client using ProxyConfig.
func WithProxy(cfg ProxyConfig) Option {
	return func(c *Client) {
		proxyURL, err := cfg.URL()
		if err != nil {
			c.configErr = err
			return
		}
		applyProxyURL(c, proxyURL)
	}
}

// WithProxyURL configures a forward proxy on the client from a raw URL string.
// Example: "https://user:pass@dc46-mum-01.algoip.in:443"
func WithProxyURL(rawURL string) Option {
	return func(c *Client) {
		trimmed := strings.TrimSpace(rawURL)
		if trimmed == "" {
			c.configErr = fmt.Errorf("%w: proxy URL cannot be empty", ErrInvalidInput)
			return
		}
		proxyURL, err := url.Parse(trimmed)
		if err != nil || proxyURL.Host == "" {
			c.configErr = fmt.Errorf("%w: invalid proxy URL %q", ErrInvalidInput, rawURL)
			return
		}
		if proxyURL.Scheme != "http" && proxyURL.Scheme != "https" {
			c.configErr = fmt.Errorf("%w: unsupported proxy scheme %q", ErrInvalidInput, proxyURL.Scheme)
			return
		}
		applyProxyURL(c, proxyURL)
	}
}

func applyProxyURL(c *Client, proxyURL *url.URL) {
	timeout := 15 * time.Second
	var tr *http.Transport

	if hc, ok := c.httpClient.(*http.Client); ok && hc != nil {
		if hc.Timeout > 0 {
			timeout = hc.Timeout
		}
		if existingTr, ok := hc.Transport.(*http.Transport); ok && existingTr != nil {
			tr = existingTr.Clone()
		}
	}

	if tr == nil {
		tr = &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	tr.Proxy = http.ProxyURL(proxyURL)

	c.httpClient = &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
}

package api

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// Client represents a Docker Registry API client
type Client struct {
	baseURL    string
	apiPrefix  string
	httpClient *http.Client
	authHeader string
}

// NewClient creates a new Docker Registry API client
func NewClient(baseURL, username, password string, insecure bool) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}

	// Remove trailing slash from baseURL if present
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Create HTTP client with optional TLS config for insecure connections
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	// Set up authentication if credentials are provided
	var authHeader string
	if username != "" && password != "" {
		auth := fmt.Sprintf("%s:%s", username, password)
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		authHeader = fmt.Sprintf("Basic %s", encodedAuth)
	}

	return &Client{
		baseURL:    baseURL,
		apiPrefix:  "/v2",
		httpClient: httpClient,
		authHeader: authHeader,
	}, nil
}

// SetAPIPrefix sets the API prefix (default is "/v2" for Docker Registry API v2)
func (c *Client) SetAPIPrefix(prefix string) {
	if prefix != "" {
		c.apiPrefix = prefix
	}
}

// buildURL builds a full URL for a Docker Registry API endpoint
func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s%s%s", c.baseURL, c.apiPrefix, path)
}

// newRequest creates a new HTTP request with proper headers
func (c *Client) newRequest(method, path string) (*http.Request, error) {
	url := c.buildURL(path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set authentication header if available
	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	}

	// Set Accept header for Docker Registry API
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	return req, nil
}

// Ping checks if the registry is reachable
func (c *Client) Ping() error {
	req, err := c.newRequest("GET", "/")
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

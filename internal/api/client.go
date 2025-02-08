package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/baudevs/yolo.baudevs.com/internal/types"
)

// Client handles API operations
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(endpoint string) *Client {
	return &Client{
		endpoint:   endpoint,
		httpClient: &http.Client{},
	}
}

// CreateCheckoutSession creates a new checkout session
func (c *Client) CreateCheckoutSession(pkg types.SubscriptionPackage) (string, error) {
	var resp struct {
		URL string `json:"url"`
	}
	if err := c.post("/checkout/create", pkg, &resp); err != nil {
		return "", err
	}
	return resp.URL, nil
}

// ActivateLicense activates a license from a session ID
func (c *Client) ActivateLicense(sessionID string) (*types.License, error) {
	var license types.License
	if err := c.post("/license/activate", map[string]string{"session_id": sessionID}, &license); err != nil {
		return nil, err
	}
	return &license, nil
}

// GetSubscriptionStatus gets the status of a subscription
func (c *Client) GetSubscriptionStatus(subscriptionID string) (bool, error) {
	var resp struct {
		Active bool `json:"active"`
	}
	if err := c.get("/subscription/"+subscriptionID+"/status", &resp); err != nil {
		return false, err
	}
	return resp.Active, nil
}

// GetCustomerSubscription gets a customer's subscription
func (c *Client) GetCustomerSubscription(customerID string) (*types.License, error) {
	var license types.License
	if err := c.get("/customer/"+customerID+"/subscription", &license); err != nil {
		return nil, err
	}
	return &license, nil
}

// GetCustomerCredits gets the current credit balance for a customer
func (c *Client) GetCustomerCredits(customerID string) (*types.Credits, error) {
	var credits types.Credits
	if err := c.get("/customer/"+customerID+"/credits", &credits); err != nil {
		return nil, err
	}
	return &credits, nil
}

// UpdateCustomerCredits updates the credit balance for a customer
func (c *Client) UpdateCustomerCredits(customerID string, credits *types.Credits) error {
	return c.post("/customer/"+customerID+"/credits", credits, nil)
}

// GetPackageCredits returns the number of credits for a price ID
func (c *Client) GetPackageCredits(priceID string) (int64, error) {
	var resp struct {
		Credits int64 `json:"credits"`
	}
	if err := c.get("/price/"+priceID+"/credits", &resp); err != nil {
		return 0, err
	}
	return resp.Credits, nil
}

// IsSubscriptionActive checks if a subscription is active
func (c *Client) IsSubscriptionActive(subscriptionID string) (bool, error) {
	return c.GetSubscriptionStatus(subscriptionID)
}

func (c *Client) get(endpoint string, resp interface{}) error {
	return c.request(http.MethodGet, endpoint, nil, resp)
}

func (c *Client) post(endpoint string, body interface{}, resp interface{}) error {
	return c.request(http.MethodPost, endpoint, body, resp)
}

func (c *Client) request(method, endpoint string, body interface{}, resp interface{}) error {
	// Build URL
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	// Marshal body if present
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", httpResp.StatusCode, respBody)
	}

	// Unmarshal response if needed
	if resp != nil {
		if err := json.Unmarshal(respBody, resp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

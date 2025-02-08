package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/types"
)

// LicenseClient handles license-related API operations
type LicenseClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewLicenseClient creates a new license client
func NewLicenseClient(baseURL string) *LicenseClient {
	return &LicenseClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// CreateCheckoutSession creates a new checkout session
func (c *LicenseClient) CreateCheckoutSession(email string, packageType string) (string, error) {
	payload := map[string]string{
		"email":   email,
		"package": packageType,
	}

	var response struct {
		SessionID string `json:"session_id"`
	}

	if err := c.post("/api/checkout/create", payload, &response); err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return response.SessionID, nil
}

// VerifyLicense validates a license key and uses credits
func (c *LicenseClient) VerifyLicense(licenseKey string, credits int) error {
	payload := map[string]interface{}{
		"license_key": licenseKey,
		"credits":     credits,
	}

	var response struct {
		Valid bool `json:"valid"`
	}

	if err := c.post("/api/license/verify", payload, &response); err != nil {
		return fmt.Errorf("failed to verify license: %w", err)
	}

	if !response.Valid {
		return fmt.Errorf("invalid license key or insufficient credits")
	}

	return nil
}

// GetLicenseAnalytics gets detailed analytics for a license
func (c *LicenseClient) GetLicenseAnalytics(licenseKey string) (*types.LicenseAnalytics, error) {
	var analytics types.LicenseAnalytics
	if err := c.get(fmt.Sprintf("/api/analytics/license/%s", licenseKey), &analytics); err != nil {
		return nil, fmt.Errorf("failed to get license analytics: %w", err)
	}
	return &analytics, nil
}

func (c *LicenseClient) post(path string, payload interface{}, response interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *LicenseClient) get(path string, response interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

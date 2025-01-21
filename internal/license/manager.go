package license

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles license operations
type Manager struct {
	config       Config
	stripeClient *StripeClient
	license      *License
}

// NewManager creates a new license manager
func NewManager(config Config) (*Manager, error) {
	stripeClient := NewStripeClient(config.StripeSecretKey)

	m := &Manager{
		config:       config,
		stripeClient: stripeClient,
	}

	// Load existing license
	if err := m.loadLicense(); err != nil {
		// Only return error if it's not a missing file
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load license: %w", err)
		}
	}

	return m, nil
}

// GetConfig returns the manager's config
func (m *Manager) GetConfig() Config {
	return m.config
}

// GetLicense returns the current license
func (m *Manager) GetLicense() *License {
	return m.license
}

// CreateCheckoutSession creates a new checkout session
func (m *Manager) CreateCheckoutSession(pkg SubscriptionPackage) (string, error) {
	return m.stripeClient.CreateCheckoutSession(pkg)
}

// ActivateSubscription activates a subscription using a session ID
func (m *Manager) ActivateSubscription(sessionID string) error {
	lic, err := m.stripeClient.ActivateSubscription(sessionID)
	if err != nil {
		return fmt.Errorf("failed to activate subscription: %w", err)
	}

	m.license = lic
	return m.saveLicense(lic)
}

// RecordAIRequest records an AI request and deducts credits
func (m *Manager) RecordAIRequest() error {
	if m.license == nil {
		return fmt.Errorf("no active license")
	}

	if !m.license.IsActive {
		return fmt.Errorf("license is not active")
	}

	// No credit deduction for unlimited plans
	if m.license.PlanType == PlanUnlimited {
		return nil
	}

	// Check if we have enough credits
	if m.license.CreditsLeft <= 0 {
		return fmt.Errorf("no credits remaining. Please purchase more credits")
	}

	// Deduct one credit
	m.license.CreditsLeft--
	return m.saveLicense(m.license)
}

// loadLicense loads the license from disk
func (m *Manager) loadLicense() error {
	path, err := getLicensePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var lic License
	if err := json.Unmarshal(data, &lic); err != nil {
		return fmt.Errorf("failed to unmarshal license: %w", err)
	}

	// Check if license has expired
	if time.Now().After(lic.ExpiresAt) {
		lic.IsActive = false
	}

	m.license = &lic
	return nil
}

// saveLicense saves the license to disk
func (m *Manager) saveLicense(lic *License) error {
	path, err := getLicensePath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create license directory: %w", err)
	}

	data, err := json.MarshalIndent(lic, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal license: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write license file: %w", err)
	}

	return nil
}

// getLicensePath returns the path to the license file
func getLicensePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".config", "yolo", "license.json"), nil
}

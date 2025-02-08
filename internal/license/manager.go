package license

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/api"
	"github.com/baudevs/yolo.baudevs.com/internal/types"
	"gopkg.in/yaml.v3"
)

// Manager handles license operations
type Manager struct {
	mu            sync.RWMutex
	license       *License
	licenseClient *api.LicenseClient
	configPath    string
}

// License represents a user's license
type License struct {
	IsActive     bool           `json:"is_active"`
	APIKey       string         `json:"api_key"`
	PlanType     types.PlanType `json:"plan_type"` // "unlimited" or "credits"
	Credits      int64          `json:"credits"`
	LastModified time.Time      `json:"last_modified"`
	CustomerID   string         `json:"customer_id,omitempty"`
	LastSyncTime time.Time      `json:"last_sync_time,omitempty"`
}

// NewManager creates a new license manager
func NewManager() (*Manager, error) {
	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	configDir = filepath.Join(configDir, "yolo")

	// Get license file path
	licensePath := filepath.Join(configDir, "settings", "license.yml")

	// If license file doesn't exist, create it with defaults
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(licensePath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create license directory: %w", err)
		}

		// Create default license
		defaultLicense := &License{
			IsActive: true,
			APIKey:   os.Getenv("OPENAI_API_KEY"), // Try to get from environment first
		}

		// Marshal to YAML
		data, err := yaml.Marshal(defaultLicense)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default license: %w", err)
		}

		// Write license file
		if err := os.WriteFile(licensePath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write license file: %w", err)
		}

		return &Manager{
			license: defaultLicense,
		}, nil
	}

	// Read existing license file
	data, err := os.ReadFile(licensePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read license file: %w", err)
	}

	// Parse license
	var license License
	if err := yaml.Unmarshal(data, &license); err != nil {
		return nil, fmt.Errorf("failed to parse license file: %w", err)
	}

	return &Manager{
		license: &license,
	}, nil
}

// GetLicense returns the current license
func (m *Manager) GetLicense() *License {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.license
}

// SaveLicense saves a new license
func (m *Manager) SaveLicense(license *License) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Set last modified time
	license.LastModified = time.Now()

	m.license = license
	return m.saveLicense()
}

// ActivateLicense activates a license with the backend
func (m *Manager) ActivateLicense(key string) error {
	// Verify license with backend
	if err := m.licenseClient.VerifyLicense(key, 0); err != nil {
		return fmt.Errorf("failed to activate license: %w", err)
	}

	// Get license analytics to get initial credit balance
	analytics, err := m.licenseClient.GetLicenseAnalytics(key)
	if err != nil {
		return fmt.Errorf("failed to get license analytics: %w", err)
	}

	// Create local license
	localLicense := &License{
		IsActive:     true,
		APIKey:       key,
		Credits:      analytics.TotalCreditsUsed,
		LastModified: time.Now(),
		LastSyncTime: time.Now(),
	}

	// Save locally
	return m.SaveLicense(localLicense)
}

// DeductCredits deducts credits from the license
func (m *Manager) DeductCredits(amount int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.license == nil || !m.license.IsActive {
		return fmt.Errorf("no active license")
	}

	// Verify with backend
	if err := m.licenseClient.VerifyLicense(m.license.APIKey, amount); err != nil {
		return fmt.Errorf("failed to verify credits: %w", err)
	}

	// Update local credit count
	m.license.Credits -= int64(amount)
	m.license.LastModified = time.Now()

	return m.saveLicense()
}

// CreateCheckoutSession creates a new checkout session for license purchase
func (m *Manager) CreateCheckoutSession(email, packageType string) (string, error) {
	sessionID, err := m.licenseClient.CreateCheckoutSession(email, packageType)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}
	return sessionID, nil
}

// loadLicense loads the license from disk
func (m *Manager) loadLicense() error {
	path, err := getLicensePath()
	if err != nil {
		return fmt.Errorf("failed to get license path: %w", err)
	}

	m.configPath = path

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		m.license = nil
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to read license file: %w", err)
	}

	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return fmt.Errorf("failed to unmarshal license: %w", err)
	}

	m.license = &license
	return nil
}

// saveLicense saves the license to disk
func (m *Manager) saveLicense() error {
	if m.configPath == "" {
		path, err := getLicensePath()
		if err != nil {
			return fmt.Errorf("failed to get license path: %w", err)
		}
		m.configPath = path
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(m.license, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal license: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write license file: %w", err)
	}

	return nil
}

// getLicensePath returns the path to the license file
func getLicensePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	return filepath.Join(configDir, "yolo", "license.json"), nil
}

// GetOpenAIKey returns the OpenAI API key from the license
func (m *Manager) GetOpenAIKey() (string, error) {
	if m.license == nil || !m.license.IsActive {
		return "", fmt.Errorf("no active license found")
	}
	if m.license.APIKey == "" {
		return "", fmt.Errorf("no API key found in license")
	}
	return m.license.APIKey, nil
}

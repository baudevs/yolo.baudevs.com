package license

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/api"
	"github.com/baudevs/yolo.baudevs.com/internal/config"
	"github.com/baudevs/yolo.baudevs.com/internal/types"
)

// Manager handles license operations
type Manager struct {
	mu         sync.RWMutex
	license    *License
	apiClient  *api.Client
	configPath string
}

// License represents a user's license
type License struct {
	IsActive      bool          `json:"is_active"`
	APIKey        string        `json:"api_key"`
	PlanType      types.PlanType `json:"plan_type"` // "unlimited" or "credits"
	Credits       int64         `json:"credits"`
	LastModified  time.Time     `json:"last_modified"`
	CustomerID    string        `json:"customer_id,omitempty"`
	LastSyncTime  time.Time     `json:"last_sync_time,omitempty"`
}

// NewManager creates a new license manager
func NewManager() (*Manager, error) {
	// Load config to get API endpoint
	cfg, err := config.LoadClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	manager := &Manager{
		apiClient: api.NewClient(cfg.APIEndpoint),
	}
	
	if err := manager.loadLicense(); err != nil {
		return nil, fmt.Errorf("failed to load license: %w", err)
	}
	
	return manager, nil
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
	// Validate license with backend
	license, err := m.apiClient.ActivateLicense(key)
	if err != nil {
		return fmt.Errorf("failed to activate license: %w", err)
	}

	// Convert API license to local license
	localLicense := &License{
		IsActive:     true,
		APIKey:       key,
		PlanType:     license.PlanType,
		Credits:      license.Credits,
		CustomerID:   license.CustomerID,
		LastModified: time.Now(),
		LastSyncTime: time.Now(),
	}

	// Save locally
	return m.SaveLicense(localLicense)
}

// SyncCredits syncs credit balance with backend
func (m *Manager) SyncCredits() error {
	m.mu.RLock()
	// Check if we need to sync
	if m.license == nil || !m.license.IsActive || m.license.CustomerID == "" {
		m.mu.RUnlock()
		return nil // Nothing to sync
	}

	// Only sync every hour
	if time.Since(m.license.LastSyncTime) < time.Hour {
		m.mu.RUnlock()
		return nil
	}
	
	customerID := m.license.CustomerID
	m.mu.RUnlock()

	// Get credits from backend
	credits, err := m.apiClient.GetCustomerCredits(customerID)
	if err != nil {
		return fmt.Errorf("failed to sync credits: %w", err)
	}

	// Update local credits
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Re-check conditions after acquiring write lock
	if m.license == nil || !m.license.IsActive || m.license.CustomerID != customerID {
		return nil
	}
	
	m.license.Credits = credits.Balance
	m.license.LastSyncTime = time.Now()
	m.license.LastModified = time.Now()

	return m.saveLicense()
}

// DeductCredits deducts credits from the license
func (m *Manager) DeductCredits(amount int) error {
	// First sync credits without holding the write lock
	if err := m.SyncCredits(); err != nil {
		return fmt.Errorf("failed to sync credits: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// No license or inactive
	if m.license == nil || !m.license.IsActive {
		return fmt.Errorf("no active license")
	}

	// Unlimited plan doesn't need credit deduction
	if m.license.PlanType == types.PlanUnlimited {
		return nil
	}

	// Convert amount to int64
	credits := int64(amount)

	// Check if we have enough credits
	if m.license.Credits < credits {
		return fmt.Errorf("insufficient credits: have %d, need %d", m.license.Credits, credits)
	}

	// Update credits in backend first
	if m.license.CustomerID != "" {
		if err := m.apiClient.UpdateCustomerCredits(m.license.CustomerID, &types.Credits{
			Balance: m.license.Credits - credits,
			UpdatedAt: time.Now(),
		}); err != nil {
			return fmt.Errorf("failed to update credits: %w", err)
		}
	}

	// Deduct credits locally
	m.license.Credits -= credits
	m.license.LastModified = time.Now()
	m.license.LastSyncTime = time.Now()

	// Save updated license
	return m.saveLicense()
}

// loadLicense loads the license from disk
func (m *Manager) loadLicense() error {
	licensePath, err := getLicensePath()
	if err != nil {
		return err
	}

	// If license doesn't exist, that's fine
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return nil
	}

	// Read license file
	data, err := os.ReadFile(licensePath)
	if err != nil {
		return fmt.Errorf("failed to read license: %w", err)
	}

	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return fmt.Errorf("failed to parse license: %w", err)
	}

	m.license = &license
	return nil
}

// saveLicense saves the license to disk
func (m *Manager) saveLicense() error {
	licensePath, err := getLicensePath()
	if err != nil {
		return err
	}

	// Create license directory if it doesn't exist
	licenseDir := filepath.Dir(licensePath)
	if err := os.MkdirAll(licenseDir, 0755); err != nil {
		return fmt.Errorf("failed to create license directory: %w", err)
	}

	// Marshal license
	data, err := json.MarshalIndent(m.license, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal license: %w", err)
	}

	// Write license file
	if err := os.WriteFile(licensePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write license: %w", err)
	}

	return nil
}

// getLicensePath returns the path to the license file
func getLicensePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".yolo", "license.json"), nil
}

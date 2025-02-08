package license

import (
	"github.com/baudevs/yolo.baudevs.com/internal/api"
	"github.com/baudevs/yolo.baudevs.com/internal/types"
)

// StripeClient handles Stripe API operations through the backend
type StripeClient struct {
	apiClient *api.Client
}

// NewStripeClient creates a new Stripe client
func NewStripeClient(apiEndpoint string) *StripeClient {
	return &StripeClient{
		apiClient: api.NewClient(apiEndpoint),
	}
}

// CreateCheckoutSession creates a new Stripe checkout session for a package
func (s *StripeClient) CreateCheckoutSession(pkg types.SubscriptionPackage) (string, error) {
	return s.apiClient.CreateCheckoutSession(pkg)
}

// ActivateSubscription activates a subscription from a checkout session
func (s *StripeClient) ActivateSubscription(sessionID string) (*types.License, error) {
	return s.apiClient.ActivateLicense(sessionID)
}

// GetSubscriptionStatus gets the status of a subscription
func (s *StripeClient) GetSubscriptionStatus(subscriptionID string) (bool, error) {
	return s.apiClient.GetSubscriptionStatus(subscriptionID)
}

// GetCustomerSubscription gets a customer's subscription
func (s *StripeClient) GetCustomerSubscription(customerID string) (*types.License, error) {
	return s.apiClient.GetCustomerSubscription(customerID)
}

// GetCustomerCredits gets the current credit balance for a customer
func (s *StripeClient) GetCustomerCredits(customerID string) (*types.Credits, error) {
	return s.apiClient.GetCustomerCredits(customerID)
}

// UpdateCustomerCredits updates the credit balance for a customer
func (s *StripeClient) UpdateCustomerCredits(customerID string, credits *types.Credits) error {
	return s.apiClient.UpdateCustomerCredits(customerID, credits)
}

// GetPackageCredits returns the number of credits for a price ID
func (s *StripeClient) GetPackageCredits(priceID string) (int64, error) {
	return s.apiClient.GetPackageCredits(priceID)
}

// IsSubscriptionActive checks if a subscription is active
func (s *StripeClient) IsSubscriptionActive(subscriptionID string) (bool, error) {
	return s.apiClient.IsSubscriptionActive(subscriptionID)
}

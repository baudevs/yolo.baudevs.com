package license

import "time"

// Config represents the license configuration
type Config struct {
	StripeSecretKey  string
	DefaultOpenAIKey string
}

// License represents a user's license
type License struct {
	CustomerID     string    `json:"customer_id"`
	SubscriptionID string    `json:"subscription_id"`
	PlanType       PlanType  `json:"plan_type"`
	LastSyncAt     time.Time `json:"last_sync_at"`
	IsActive       bool      `json:"is_active"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreditsLeft    int64     `json:"credits_left"`
	APIKey         string    `json:"api_key"`
}

// PlanType represents the type of subscription plan
type PlanType string

const (
	PlanStarter    PlanType = "starter"
	PlanPro        PlanType = "pro"
	PlanTeam       PlanType = "team"
	PlanEnterprise PlanType = "enterprise"
	PlanUnlimited  PlanType = "unlimited"
)

// SubscriptionPackage represents a subscription package
type SubscriptionPackage struct {
	Type        PlanType `json:"type"`
	Name        string   `json:"name"`
	PriceID     string   `json:"price_id"`   // Stripe Price ID
	ProductID   string   `json:"product_id"` // Stripe Product ID
	Description string   `json:"description"`
	BasePrice   float64  `json:"base_price"` // Base price per month
	Credits     int64    `json:"credits,omitempty"`
}

// Available packages
var (
	StarterPack = SubscriptionPackage{
		Type:        PlanStarter,
		Name:        "Starter Pack",
		Description: "Perfect for personal projects. 5,000 credits.",
		BasePrice:   5.00,
		Credits:     5_000,
	}

	ProPack = SubscriptionPackage{
		Type:        PlanPro,
		Name:        "Pro Pack",
		Description: "For power users. 25,000 credits.",
		BasePrice:   20.00,
		Credits:     25_000,
	}

	TeamPack = SubscriptionPackage{
		Type:        PlanTeam,
		Name:        "Team Pack",
		Description: "Great for small teams. 100,000 credits.",
		BasePrice:   60.00,
		Credits:     100_000,
	}

	EnterprisePack = SubscriptionPackage{
		Type:        PlanEnterprise,
		Name:        "Enterprise Pack",
		Description: "For large teams. 500,000 credits.",
		BasePrice:   200.00,
		Credits:     500_000,
	}

	UnlimitedPack = SubscriptionPackage{
		Type:        PlanUnlimited,
		Name:        "Unlimited Monthly",
		Description: "Unlimited credits, monthly subscription.",
		BasePrice:   99.00,
		Credits:     -1, // Unlimited
	}
)

// GetPersonalPackages returns all personal packages
func GetPersonalPackages() []SubscriptionPackage {
	return []SubscriptionPackage{
		StarterPack,
		ProPack,
	}
}

// GetBusinessPackages returns all business packages
func GetBusinessPackages() []SubscriptionPackage {
	return []SubscriptionPackage{
		TeamPack,
		EnterprisePack,
		UnlimitedPack,
	}
}

// GetPackageByType returns a package by its type
func GetPackageByType(planType PlanType) *SubscriptionPackage {
	switch planType {
	case PlanStarter:
		return &StarterPack
	case PlanPro:
		return &ProPack
	case PlanTeam:
		return &TeamPack
	case PlanEnterprise:
		return &EnterprisePack
	case PlanUnlimited:
		return &UnlimitedPack
	default:
		return nil
	}
}

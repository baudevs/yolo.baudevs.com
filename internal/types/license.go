package types

import "time"

// PlanType represents a subscription plan type
type PlanType string

const (
	PlanStarter    PlanType = "starter"
	PlanPro        PlanType = "pro"
	PlanTeam       PlanType = "team"
	PlanEnterprise PlanType = "enterprise"
	PlanUnlimited  PlanType = "unlimited"
)

// SubscriptionPackage represents a subscription package type
type SubscriptionPackage string

const (
	PackageStarter    SubscriptionPackage = "starter"
	PackagePro        SubscriptionPackage = "pro"
	PackageTeam       SubscriptionPackage = "team"
	PackageEnterprise SubscriptionPackage = "enterprise"
	PackageUnlimited  SubscriptionPackage = "unlimited"
)

// License represents a YOLO license
type License struct {
	CustomerID   string    `json:"customer_id"`
	PlanType     PlanType  `json:"plan_type"`
	IsActive     bool      `json:"is_active"`
	ActivatedAt  time.Time `json:"activated_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Credits      int64     `json:"credits"`
	SessionID    string    `json:"session_id"`
	LastModified time.Time `json:"last_modified"`
	APIKey       string    `json:"api_key"`
}

// RateLimit represents rate limiting configuration
type RateLimit struct {
	MaxQueries int   `json:"max_queries"` // Maximum number of queries
	Period     int64 `json:"period"`      // Period in seconds
}

// Credits represents the credit metadata stored in Stripe
type Credits struct {
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Available packages
var (
	StarterPackage = SubscriptionPackage("starter")
	ProPackage     = SubscriptionPackage("pro")
	TeamPackage    = SubscriptionPackage("team")
)

// GetPackage returns a package by type
func GetPackage(typ PlanType) *SubscriptionPackage {
	switch typ {
	case PlanStarter:
		return &StarterPackage
	case PlanPro:
		return &ProPackage
	case PlanTeam:
		return &TeamPackage
	default:
		return nil
	}
}

// AllPackages returns all available packages
func AllPackages() []SubscriptionPackage {
	return []SubscriptionPackage{
		StarterPackage,
		ProPackage,
		TeamPackage,
	}
}

// LicenseAnalytics represents analytics data for a license
type LicenseAnalytics struct {
	UsageStats       []UsageStat `json:"usage_stats"`
	ValidationCount  int         `json:"validation_count"`
	LastUsed         time.Time   `json:"last_used"`
	TotalCreditsUsed int64       `json:"total_credits_used"`
}

// UsageStat represents a single usage statistic
type UsageStat struct {
	Timestamp time.Time `json:"timestamp"`
	Credits   int64     `json:"credits"`
	Action    string    `json:"action"`
}

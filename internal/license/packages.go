package license

import "time"

// PackageType represents the type of license package
type PackageType string

const (
	PackageStarter    PackageType = "starter"
	PackagePro        PackageType = "pro"
	PackageTeam       PackageType = "team"
	PackageEnterprise PackageType = "enterprise"
	PackageUnlimited  PackageType = "unlimited"
)

// CreditPackage represents a one-time credit package configuration
type CreditPackage struct {
	Type        PackageType    `json:"type"`
	Name        string         `json:"name"`
	Credits     int64          `json:"credits"`
	Price       float64        `json:"price"`      // in USD
	Duration    *time.Duration `json:"duration"`   // nil for one-time packages
	RateLimit   *RateLimit     `json:"rate_limit"` // nil for unlimited rate
}

// Available credit packages
var (
	StarterPackage = CreditPackage{
		Type:    PackageStarter,
		Name:    "Starter Pack",
		Credits: 5_000,
		Price:   5.00,
	}

	ProPackage = CreditPackage{
		Type:    PackagePro,
		Name:    "Pro Pack",
		Credits: 25_000,
		Price:   20.00,
	}

	TeamPackage = CreditPackage{
		Type:    PackageTeam,
		Name:    "Team Pack",
		Credits: 100_000,
		Price:   75.00,
		RateLimit: &RateLimit{
			MaxQueries: 1000,
			Period:     3600, // 1 hour
		},
	}
)

// GetPackage returns a package by type
func GetPackage(typ PackageType) *CreditPackage {
	switch typ {
	case PackageStarter:
		return &StarterPackage
	case PackagePro:
		return &ProPackage
	case PackageTeam:
		return &TeamPackage
	default:
		return nil
	}
}

// AllPackages returns all available packages
func AllPackages() []CreditPackage {
	return []CreditPackage{
		StarterPackage,
		ProPackage,
		TeamPackage,
	}
}

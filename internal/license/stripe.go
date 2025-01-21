package license

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/usagerecord"
)

const (
	MeterAIRequests = "ai_requests"
)

// StripeClient handles Stripe API operations
type StripeClient struct {
	secretKey string
}

// NewStripeClient creates a new Stripe client
func NewStripeClient(secretKey string) *StripeClient {
	stripe.Key = secretKey
	return &StripeClient{
		secretKey: secretKey,
	}
}

// RecordAIRequest records an AI request usage event
func (s *StripeClient) RecordAIRequest(subscriptionItemID string) error {
	params := &stripe.UsageRecordParams{
		SubscriptionItem: stripe.String(subscriptionItemID),
		Quantity:        stripe.Int64(1),
		Timestamp:       stripe.Int64(time.Now().Unix()),
		Action:          stripe.String(string(stripe.UsageRecordActionIncrement)),
	}

	_, err := usagerecord.New(params)
	if err != nil {
		return fmt.Errorf("failed to record AI request: %w", err)
	}

	return nil
}

// CreateCheckoutSession creates a new Stripe checkout session for a package
func (s *StripeClient) CreateCheckoutSession(pkg SubscriptionPackage) (string, error) {
	// Create price ID if not set
	if pkg.PriceID == "" {
		price, err := s.createPrice(pkg)
		if err != nil {
			return "", fmt.Errorf("failed to create price: %w", err)
		}
		pkg.PriceID = price.ID
	}

	// Create checkout session
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(pkg.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("https://yolo.baudevs.com/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://yolo.baudevs.com/cancel"),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess.URL, nil
}

// ActivateSubscription activates a subscription from a checkout session
func (s *StripeClient) ActivateSubscription(sessionID string) (*License, error) {
	// Get checkout session
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkout session: %w", err)
	}

	// Get subscription
	sub, err := subscription.Get(sess.Subscription.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Create license
	lic := &License{
		CustomerID:     sess.Customer.ID,
		SubscriptionID: sess.Subscription.ID,
		PlanType:       getPlanTypeFromPrice(sub.Items.Data[0].Price.ID),
		LastSyncAt:     time.Now(),
		IsActive:       true,
		ExpiresAt:      time.Unix(sub.CurrentPeriodEnd, 0),
		CreditsLeft:    getCreditsFromPlan(getPlanTypeFromPrice(sub.Items.Data[0].Price.ID)),
	}

	return lic, nil
}

// GetSubscriptionStatus gets the status of a subscription
func (s *StripeClient) GetSubscriptionStatus(subscriptionID string) (bool, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub.Status == stripe.SubscriptionStatusActive, nil
}

// GetCustomerSubscription gets a customer's subscription
func (s *StripeClient) GetCustomerSubscription(customerID string) (*stripe.Subscription, error) {
	c, err := s.GetCustomer(customerID)
	if err != nil {
		return nil, err
	}

	if len(c.Subscriptions.Data) == 0 {
		return nil, fmt.Errorf("no active subscription found")
	}

	return c.Subscriptions.Data[0], nil
}

// GetCustomer gets a customer from Stripe
func (s *StripeClient) GetCustomer(customerID string) (*stripe.Customer, error) {
	return customer.Get(customerID, nil)
}

// GetCustomerCredits gets the current credit balance for a customer
func (s *StripeClient) GetCustomerCredits(customerID string) (*Credits, error) {
	c, err := customer.Get(customerID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	var credits Credits
	if c.Metadata["credits"] != "" {
		if err := json.Unmarshal([]byte(c.Metadata["credits"]), &credits); err != nil {
			return nil, fmt.Errorf("failed to unmarshal credits: %w", err)
		}
	}

	return &credits, nil
}

// UpdateCustomerCredits updates the credit balance for a customer
func (s *StripeClient) UpdateCustomerCredits(customerID string, credits *Credits) error {
	creditsJSON, err := json.Marshal(credits)
	if err != nil {
		return fmt.Errorf("failed to marshal credits: %w", err)
	}

	params := &stripe.CustomerParams{
		Metadata: map[string]string{
			"credits": string(creditsJSON),
		},
	}

	_, err = customer.Update(customerID, params)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	return nil
}

// GetPackageCredits returns the number of credits for a price ID
func (s *StripeClient) GetPackageCredits(priceID string) (int64, error) {
	p, err := price.Get(priceID, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get price: %w", err)
	}

	credits, ok := p.Metadata["credits"]
	if !ok {
		return 0, fmt.Errorf("no credits defined for price %s", priceID)
	}

	var amount int64
	if err := json.Unmarshal([]byte(credits), &amount); err != nil {
		return 0, fmt.Errorf("failed to unmarshal credits: %w", err)
	}

	return amount, nil
}

// IsSubscriptionActive checks if a subscription is active
func (s *StripeClient) IsSubscriptionActive(subscriptionID string) (bool, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub.Status == stripe.SubscriptionStatusActive, nil
}

// getPlanTypeFromPrice returns the plan type for a given price ID
func getPlanTypeFromPrice(priceID string) PlanType {
	switch priceID {
	case StarterPack.PriceID:
		return PlanStarter
	case ProPack.PriceID:
		return PlanPro
	case TeamPack.PriceID:
		return PlanTeam
	case EnterprisePack.PriceID:
		return PlanEnterprise
	case UnlimitedPack.PriceID:
		return PlanUnlimited
	default:
		return PlanStarter
	}
}

// getCreditsFromPlan returns the number of credits for a given plan type
func getCreditsFromPlan(planType PlanType) int64 {
	switch planType {
	case PlanStarter:
		return StarterPack.Credits
	case PlanPro:
		return ProPack.Credits
	case PlanTeam:
		return TeamPack.Credits
	case PlanEnterprise:
		return EnterprisePack.Credits
	case PlanUnlimited:
		return -1 // Unlimited
	default:
		return StarterPack.Credits
	}
}

// createPrice creates a new Stripe price for a package
func (s *StripeClient) createPrice(pkg SubscriptionPackage) (*stripe.Price, error) {
	// First create a product
	prod, err := product.New(&stripe.ProductParams{
		Name:        stripe.String(pkg.Name),
		Description: stripe.String(pkg.Description),
		Metadata: map[string]string{
			"plan_type": string(pkg.Type),
			"credits":   fmt.Sprintf("%d", pkg.Credits),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Then create a price for the product
	params := &stripe.PriceParams{
		Currency:   stripe.String("usd"),
		UnitAmount: stripe.Int64(int64(pkg.BasePrice * 100)), // Convert to cents
		Product:    stripe.String(prod.ID),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
	}

	return price.New(params)
}

// Credits represents the credit metadata stored in Stripe
type Credits struct {
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

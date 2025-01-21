// Package api provides webhook handlers
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type WebhookResponse struct {
	Error string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func logError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func logInfo(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Handler handles Stripe webhook events
func Handler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Stripe-Signature")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, WebhookResponse{Error: "method not allowed"})
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logError("Failed to read request body: %v", err)
		writeJSON(w, http.StatusBadRequest, WebhookResponse{Error: "failed to read request body"})
		return
	}

	// Get Stripe signature
	stripeSignature := r.Header.Get("Stripe-Signature")
	if stripeSignature == "" {
		logError("Missing Stripe signature")
		writeJSON(w, http.StatusBadRequest, WebhookResponse{Error: "missing stripe signature"})
		return
	}

	// Initialize Stripe
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		logError("STRIPE_SECRET_KEY not set")
		writeJSON(w, http.StatusInternalServerError, WebhookResponse{Error: "stripe key not configured"})
		return
	}
	stripe.Key = stripeKey

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		logError("STRIPE_WEBHOOK_SECRET not set")
		writeJSON(w, http.StatusInternalServerError, WebhookResponse{Error: "webhook secret not configured"})
		return
	}

	logInfo("Webhook signature: %s", stripeSignature)

	// Verify signature and construct event
	event, err := webhook.ConstructEvent(body, stripeSignature, endpointSecret)
	if err != nil {
		logError("Failed to verify webhook signature: %v", err)
		writeJSON(w, http.StatusBadRequest, WebhookResponse{Error: "failed to verify webhook signature"})
		return
	}

	// Handle event
	switch event.Type {
	case "charge.succeeded":
		handleChargeSucceeded(event)
	case "customer.subscription.created":
		handleSubscriptionCreated(event)
	case "customer.subscription.updated":
		handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		handleSubscriptionDeleted(event)
	case "invoice.paid":
		handleInvoicePaid(event)
	case "invoice.payment_failed":
		handleInvoicePaymentFailed(event)
	case "payment_intent.succeeded":
		handlePaymentIntentSucceeded(event)
	default:
		logInfo("Unhandled event type: %s", event.Type)
	}

	writeJSON(w, http.StatusOK, WebhookResponse{})
}

func handleChargeSucceeded(event stripe.Event) {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		logError("Failed to unmarshal charge: %v", err)
		return
	}

	logInfo("Charge succeeded: %s", charge.ID)
}

func handleSubscriptionCreated(event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		logError("Failed to unmarshal subscription: %v", err)
		return
	}

	// Initialize license manager
	manager, err := license.NewManager(license.Config{
		StripeSecretKey: stripe.Key,
	})
	if err != nil {
		logError("Failed to initialize license manager: %v", err)
		return
	}

	// Activate subscription
	if err := manager.ActivateSubscription(subscription.ID); err != nil {
		logError("Failed to activate subscription: %v", err)
		return
	}

	logInfo("Subscription created and activated: %s", subscription.ID)
}

func handleSubscriptionUpdated(event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		logError("Failed to unmarshal subscription: %v", err)
		return
	}

	// Initialize license manager
	manager, err := license.NewManager(license.Config{
		StripeSecretKey: stripe.Key,
	})
	if err != nil {
		logError("Failed to initialize license manager: %v", err)
		return
	}

	// Update subscription
	if err := manager.ActivateSubscription(subscription.ID); err != nil {
		logError("Failed to update subscription: %v", err)
		return
	}

	logInfo("Subscription updated: %s", subscription.ID)
}

func handleSubscriptionDeleted(event stripe.Event) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		logError("Failed to unmarshal subscription: %v", err)
		return
	}

	logInfo("Subscription deleted: %s", subscription.ID)
}

func handleInvoicePaid(event stripe.Event) {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		logError("Failed to unmarshal invoice: %v", err)
		return
	}

	logInfo("Invoice paid: %s", invoice.ID)
}

func handleInvoicePaymentFailed(event stripe.Event) {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		logError("Failed to unmarshal invoice: %v", err)
		return
	}

	logInfo("Invoice payment failed: %s", invoice.ID)
}

func handlePaymentIntentSucceeded(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		logError("Failed to unmarshal payment intent: %v", err)
		return
	}

	logInfo("Payment intent succeeded: %s", paymentIntent.ID)
}

package payments

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

type stripeClient struct {
	webhookSecret string
}

func NewStripeClient(key string) StripeClient {
	stripe.Key = key
	return &stripeClient{}
}

// SetWebhookSecret sets the webhook secret for signature verification.
func (s *stripeClient) SetWebhookSecret(secret string) {
	s.webhookSecret = secret
}

func (s *stripeClient) CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}

	// Add metadata
	for k, v := range metadata {
		params.AddMetadata(k, v)
	}

	// Automatic payment methods for easier integration
	params.AutomaticPaymentMethods = &stripe.PaymentIntentAutomaticPaymentMethodsParams{
		Enabled: stripe.Bool(true),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}

	return pi.ID, nil
}

func (s *stripeClient) VerifyWebhookSignature(payload []byte, sigHeader string) (bool, error) {
	if s.webhookSecret == "" {
		// In test/dev mode without webhook secret, skip verification
		return true, nil
	}

	_, err := webhook.ConstructEvent(payload, sigHeader, s.webhookSecret)
	if err != nil {
		return false, err
	}

	return true, nil
}

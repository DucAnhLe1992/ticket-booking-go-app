package payments

// Stripe client stub. Replace with a real stripe-go wrapper.

type stripeClientStub struct {
	key string
}

func NewStripeClient(key string) StripeClient {
	return &stripeClientStub{key: key}
}

func (s *stripeClientStub) CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (string, error) {
	// In real implementation: use stripe-go to create a PaymentIntent and return its ID.
	return "pi_stub_123", nil
}

func (s *stripeClientStub) VerifyWebhookSignature(payload []byte, sigHeader string) (bool, error) {
	// In real implementation: use stripe.Webhook.ConstructEvent or similar to verify.
	return true, nil
}

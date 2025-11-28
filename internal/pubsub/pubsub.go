package pubsub

import "context"

// Publisher is a minimal pub/sub publisher interface used by services.
type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Close() error
}

// Subscriber is a minimal subscriber interface.
type Subscriber interface {
	Subscribe(subject string, handler func(msg []byte)) error
	Close() error
}

// NewNoopPublisher returns a publisher that does nothing (useful for dev).
func NewNoopPublisher() Publisher { return noop{} }

type noop struct{}

func (n noop) Publish(ctx context.Context, subject string, data []byte) error { return nil }
func (n noop) Close() error                                                  { return nil }
func (n noop) Subscribe(subject string, handler func(msg []byte)) error     { return nil }

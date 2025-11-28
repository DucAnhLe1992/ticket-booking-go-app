package pubsub

import (
	"context"
	"time"

	nats "github.com/nats-io/nats.go"
)

// NATSClient is a thin wrapper around nats.Conn implementing Publisher and Subscriber.
type NATSClient struct {
	conn *nats.Conn
}

// NewNATS connects to a NATS server at the given URL and returns a Client.
// Example URL: nats://localhost:4222
func NewNATS(url string, opts ...nats.Option) (Client, error) {
	// sensible defaults: name, reconnect, timeouts
	base := []nats.Option{
		nats.Name("ticket-booking-go-app"),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1),
		nats.Timeout(5 * time.Second),
		nats.PingInterval(20 * time.Second),
		nats.FlusherTimeout(2 * time.Second),
	}
	base = append(base, opts...)
	conn, err := nats.Connect(url, base...)
	if err != nil {
		return nil, err
	}
	return &NATSClient{conn: conn}, nil
}

func (c *NATSClient) Publish(ctx context.Context, subject string, data []byte) error {
	// Publish with context-aware deadline if present
	if deadline, ok := ctx.Deadline(); ok {
		// Use RequestWithContext style semantics by using Publish and Flush within deadline
		c.conn.Publish(subject, data)
		t := time.Until(deadline)
		if t <= 0 {
			t = time.Millisecond
		}
		return c.conn.FlushTimeout(t)
	}
	return c.conn.Publish(subject, data)
}

func (c *NATSClient) Subscribe(subject string, handler func(msg []byte)) error {
	_, err := c.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	})
	return err
}

func (c *NATSClient) Close() error {
	if c.conn != nil && !c.conn.IsClosed() {
		c.conn.Drain()
		c.conn.Close()
	}
	return nil
}

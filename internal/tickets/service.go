package tickets

import (
	"context"
	"encoding/json"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

type Service struct {
	repo Repository
	pub  pubsub.Publisher
}

func NewService(repo Repository, pub pubsub.Publisher) *Service {
	return &Service{repo: repo, pub: pub}
}

func (s *Service) Create(ctx context.Context, title string, price int64, userID string) (*Ticket, error) {
	t, err := s.repo.Create(ctx, title, price, userID)
	if err != nil {
		return nil, err
	}
	if s.pub != nil {
		evt := events.TicketCreatedData{ID: t.ID, Title: t.Title, Price: t.Price, UserID: t.UserID, Version: t.Version}
		b, _ := json.Marshal(evt)
		_ = s.pub.Publish(ctx, string(events.SubjectTicketCreated), b)
	}
	return t, nil
}

func (s *Service) Update(ctx context.Context, id string, version int, title string, price int64, userID string) (*Ticket, error) {
	t, err := s.repo.UpdateWithVersion(ctx, id, version, title, price, userID)
	if err != nil {
		return nil, err
	}
	if s.pub != nil {
		evt := events.TicketUpdatedData{ID: t.ID, Title: t.Title, Price: t.Price, UserID: t.UserID, OrderID: t.OrderID, Version: t.Version}
		b, _ := json.Marshal(evt)
		_ = s.pub.Publish(ctx, string(events.SubjectTicketUpdated), b)
	}
	return t, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Ticket, error) { return s.repo.Get(ctx, id) }
func (s *Service) List(ctx context.Context) ([]*Ticket, error)         { return s.repo.List(ctx) }

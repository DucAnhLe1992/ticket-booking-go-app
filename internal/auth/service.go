package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"encoding/json"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      UserRepository
	validator *validator.Validate
	pub       pubsub.Publisher
}

func NewService(repo UserRepository, pub pubsub.Publisher) *Service {
	return &Service{repo: repo, validator: validator.New(), pub: pub}
}

type SignupInput struct {
	Email    string `json:"email" validate:"required,email,max=200"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type SigninInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *Service) Signup(ctx context.Context, in SignupInput) (*User, string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, "", err
	}
	existing, err := s.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("email in use")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	u, err := s.repo.Create(ctx, in.Email, string(hash))
	if err != nil {
		return nil, "", err
	}
	token, err := issueJWT(u.ID, u.Email)
	if err != nil {
		return nil, "", err
	}
	// emit user:created (optional)
	if s.pub != nil {
		evt := events.UserCreatedData{ID: u.ID, Email: u.Email}
		b, _ := json.Marshal(evt)
		_ = s.pub.Publish(ctx, string(events.SubjectUserCreated), b)
	}
	return u, token, nil
}

func (s *Service) Signin(ctx context.Context, in SigninInput) (*User, string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, "", err
	}
	u, err := s.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		return nil, "", err
	}
	if u == nil {
		return nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := issueJWT(u.ID, u.Email)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func issueJWT(id, email string) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "dev-secret-key" // dev fallback; use real secret in prod
	}
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

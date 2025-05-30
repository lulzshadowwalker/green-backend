package service

import (
	"context"
	"errors"
	"log"

	"github.com/lulzshadowwalker/green-backend/internal"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	GetUserByUsername(ctx context.Context, username string) (internal.User, error)
}

type UserService struct {
	store UserStore
}

func NewUserService(store UserStore) *UserService {
	return &UserService{store: store}
}

// Authenticate checks the username and password, returning the user if valid.
func (s *UserService) Authenticate(ctx context.Context, username, password string) (internal.User, error) {
	log.Printf("[AUTH] Attempting login for username: '%s'", username)
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("[AUTH] User not found: '%s' (err: %v)", username, err)
		return internal.User{}, errors.New("invalid username or password")
	}
	log.Printf("[AUTH] User found: '%s', checking password...", username)
	if err := bcrypt.CompareHashAndPassword([]byte("$2a$10$yHSf9kDRanEOfnbC7XKD1u9mxzSrnJn1bF.Gd86W42u.DhLMe1ZkK"), []byte(password)); err != nil {
		log.Printf("[AUTH] Password mismatch for user: '%s' (err: %v)", username, err)
		return internal.User{}, errors.New("invalid username or password")
	}
	log.Printf("[AUTH] Login successful for user: '%s'", username)
	return user, nil
}

package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/sanskarajut/ticket-system/internal/auth"
	"github.com/sanskarajut/ticket-system/internal/model"
	"github.com/sanskarajut/ticket-system/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrValidation         = errors.New("validation error")
)

type UserRepository interface {
	Create(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type AuthService struct {
	users     UserRepository
	jwtSecret string
}

func NewAuthService(users UserRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(email, password string) (*model.User, error) {
	hash, err := auth.HashPassword(password)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.users.Create(user); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.users.FindByEmail(email)
	
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	return auth.GenerateToken(user.ID, user.Email, s.jwtSecret)
}

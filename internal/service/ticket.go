package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/sanskarajut/ticket-system/internal/model"
	"github.com/sanskarajut/ticket-system/internal/repository"
)

type TicketRepository interface {
	Create(t *model.Ticket) error
	FindByID(id string) (*model.Ticket, error)
	FindByUserID(userID string) ([]*model.Ticket, error)
	UpdateStatus(id string, status model.TicketStatus, updatedAt interface{}) error
}

type TicketService struct {
	tickets TicketRepository
}

func NewTicketService(tickets TicketRepository) *TicketService {
	return &TicketService{tickets: tickets}
}

func (s *TicketService) Create(userID, title, description string) (*model.Ticket, error) {
	if title == "" {
		return nil, ErrValidation
	}

	now := time.Now().UTC()

	t := &model.Ticket{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Desc:      description,
		Status:    model.StatusOpen,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.tickets.Create(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TicketService) ListForUser(userID string) ([]*model.Ticket, error) {
	return s.tickets.FindByUserID(userID)
}

func (s *TicketService) GetForUser(userID, ticketID string) (*model.Ticket, error) {
	t, err := s.tickets.FindByID(ticketID)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if t.UserID != userID {
		return nil, ErrForbidden
	}

	return t, nil
}

func (s *TicketService) UpdateStatus(userID, ticketID string, next model.TicketStatus) (*model.Ticket, error) {
	t, err := s.GetForUser(userID, ticketID)

	if err != nil {
		return nil, err
	}

	if !t.Status.ValidTransition(next) {
		return nil, ErrInvalidTransition
	}

	now := time.Now().UTC()
	if err := s.tickets.UpdateStatus(ticketID, next, now); err != nil {
		return nil, err
	}

	t.Status = next
	t.UpdatedAt = now
	
	return t, nil
}

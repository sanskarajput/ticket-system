package model

import "time"

type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusClosed     TicketStatus = "closed"
)

// ValidTransition returns true if moving from current -> next is allowed.
func (s TicketStatus) ValidTransition(next TicketStatus) bool {

	switch s {
	case StatusOpen:
		return next == StatusInProgress
	case StatusInProgress:
		return next == StatusClosed
	default:
		return false
	}
}

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}




type Ticket struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Title     string       `json:"title"`
	Desc      string       `json:"description"`
	Status    TicketStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}        







package repository

import (
	"database/sql"
	"errors"

	"github.com/sanskarajut/ticket-system/internal/model"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(t *model.Ticket) error {
	_, err := r.db.Exec(
		`INSERT INTO tickets (id, user_id, title, description, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.UserID, t.Title, t.Desc, t.Status, t.CreatedAt, t.UpdatedAt,
	)

	return err
}

func (r *TicketRepository) FindByID(id string) (*model.Ticket, error) {
	t := &model.Ticket{}

	err := r.db.QueryRow(
		`SELECT id, user_id, title, description, status, created_at, updated_at
		 FROM tickets WHERE id = ?`, id,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Desc, &t.Status, &t.CreatedAt, &t.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return t, err
}

func (r *TicketRepository) FindByUserID(userID string) ([]*model.Ticket, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, description, status, created_at, updated_at
		 FROM tickets WHERE user_id = ? ORDER BY created_at DESC`, userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*model.Ticket

	for rows.Next() {
		t := &model.Ticket{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Desc, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (r *TicketRepository) UpdateStatus(id string, status model.TicketStatus, updatedAt interface{}) error {
	_, err := r.db.Exec(
		`UPDATE tickets SET status = ?, updated_at = ? WHERE id = ?`,
		status, updatedAt, id,
	)

	return err
}



             
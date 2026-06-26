package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sanskarajut/ticket-system/internal/middleware"
	"github.com/sanskarajut/ticket-system/internal/model"
	"github.com/sanskarajut/ticket-system/internal/service"
)

type TicketHandler struct {
	svc *service.TicketService
}

func newTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

type createTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	ticket, err := h.svc.Create(userID, req.Title, req.Description)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	tickets, err := h.svc.ListForUser(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if tickets == nil {
		tickets = []*model.Ticket{}
	}

	writeJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	ticketID := r.PathValue("id")

	ticket, err := h.svc.GetForUser(userID, ticketID)

	if err != nil {
		switch err {
		case service.ErrNotFound:
			writeError(w, http.StatusNotFound, "ticket not found")
		case service.ErrForbidden:
			writeError(w, http.StatusForbidden, "access denied")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
             
	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	ticketID := r.PathValue("id")

	var req updateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	next := model.TicketStatus(req.Status)

	switch next {
	case model.StatusOpen, model.StatusInProgress, model.StatusClosed:
	default:              
		writeError(w, http.StatusBadRequest, "invalid status; must be open, in_progress, or closed")
		return
	}

	ticket, err := h.svc.UpdateStatus(userID, ticketID, next)

	if err != nil {
		switch err {
		case service.ErrNotFound:
			writeError(w, http.StatusNotFound, "ticket not found")
		case service.ErrForbidden:
			writeError(w, http.StatusForbidden, "access denied")
		case service.ErrInvalidTransition:
			writeError(w, http.StatusUnprocessableEntity, "invalid status transition")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}













 
 


 

 
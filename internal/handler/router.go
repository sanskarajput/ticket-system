package handler

import (
	"net/http"

	"github.com/sanskarajut/ticket-system/internal/middleware"
	"github.com/sanskarajut/ticket-system/internal/service"
)

type Handler struct {
	auth   *AuthHandler
	ticket *TicketHandler
}

func New(authSvc *service.AuthService, ticketSvc *service.TicketService) *Handler {
	return &Handler{
		auth:   newAuthHandler(authSvc),
		ticket: newTicketHandler(ticketSvc),
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux, authMw *middleware.AuthMiddleware) {
	// Public routes
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	router.HandleFunc("POST /auth/register", h.auth.Register)
	router.HandleFunc("POST /auth/login", h.auth.Login)

	// Protected routes
	router.Handle("POST /tickets", authMw.Require(http.HandlerFunc(h.ticket.Create)))
	router.Handle("GET /tickets", authMw.Require(http.HandlerFunc(h.ticket.List)))
	router.Handle("GET /tickets/{id}", authMw.Require(http.HandlerFunc(h.ticket.Get)))
	router.Handle("PATCH /tickets/{id}/status", authMw.Require(http.HandlerFunc(h.ticket.UpdateStatus)))
}

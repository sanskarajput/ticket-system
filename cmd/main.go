package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sanskarajut/ticket-system/internal/config"
	"github.com/sanskarajut/ticket-system/internal/db"
	"github.com/sanskarajut/ticket-system/internal/handler"
	"github.com/sanskarajut/ticket-system/internal/middleware"
	"github.com/sanskarajut/ticket-system/internal/repository"
	"github.com/sanskarajut/ticket-system/internal/service"
)

func main() {
	// load config
	_ = godotenv.Load()
	cfg := config.Load()

	// database setup
	database, err := db.Init(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	ticketRepo := repository.NewTicketRepository(database)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	ticketSvc := service.NewTicketService(ticketRepo)

	authMw := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// setup router
	router := http.NewServeMux()

	h := handler.New(authSvc, ticketSvc)
	h.RegisterRoutes(router, authMw)

	// setup server
	log.Printf("server listening on http://localhost:%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

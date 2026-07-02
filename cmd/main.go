package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/sanskarajut/ticket-system/internal/config"
	"github.com/sanskarajut/ticket-system/internal/db"
	"github.com/sanskarajut/ticket-system/internal/handler"
	"github.com/sanskarajut/ticket-system/internal/middleware"
	"github.com/sanskarajut/ticket-system/internal/repository"
	"github.com/sanskarajut/ticket-system/internal/service"
)

// Target URL for the internal health check
const apiURL = "https://sanskarajput-s-ticket-system.onrender.com/health"

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

	// --- START OF BACKGROUND HEALTH WORKER ---
	go func() {
		// 1. Run immediately on startup
		pingHealth()

		// 2. Repeat every 5 seconds
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			pingHealth()
		}
	}()
	// --- END OF BACKGROUND HEALTH WORKER ---

	// setup server
	log.Printf("server listening on http://localhost:%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// pingHealth mimics your JavaScript function using Go's http client
func pingHealth() {
	// Use a client timeout so it doesn't hang indefinitely if Render goes down
	client := &http.Client{Timeout: 4 * time.Second}

	resp, err := client.Get(apiURL)
	timeStr := time.Now().Format("15:04:05")

	if err != nil {
		log.Printf("[%s] Health check failed: %v", timeStr, err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON data dynamically like the JS client
	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("[%s] Failed to parse JSON response: %v", timeStr, err)
		return
	}

	log.Printf("[%s] %+v", timeStr, data)
}

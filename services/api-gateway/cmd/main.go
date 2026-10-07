package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/maxluvds/booking-platform/pkg/config"
	"github.com/maxluvds/booking-platform/pkg/logger"
	"github.com/maxluvds/booking-platform/services/api-gateway/internal/client"
	"github.com/maxluvds/booking-platform/services/api-gateway/internal/handler"
	"github.com/maxluvds/booking-platform/services/api-gateway/internal/middleware"
)

func main() {
	cfg, err := config.Load("services/api-gateway/config/config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	log := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)

	log.Info("API Gateway starting...")

	eventClient, err := client.NewEventClient(
		cfg.Services.EventService.URL,
		cfg.Services.EventService.Timeout,
	)
	if err != nil {
		log.Error("Failed to create event client", "error", err)
		return
	}
	defer eventClient.Close()
	log.Info("Connected to Event Service", "url", cfg.Services.EventService.URL)

	userClient, err := client.NewUserClient(
		cfg.Services.UserService.URL,
		cfg.Services.UserService.Timeout,
	)
	if err != nil {
		log.Error("Failed to create user client", "error", err)
		return
	}
	defer userClient.Close()
	log.Info("Connected to User Service", "url", cfg.Services.UserService.URL)

	eventHandler := handler.NewEventHandler(eventClient)
	authHandler := handler.NewAuthHandler(userClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	authMiddleware := middleware.JWTAuth(userClient, log)
	mux.Handle("GET /api/v1/auth/me", authMiddleware(http.HandlerFunc(authHandler.GetMe)))

	mux.Handle("POST /api/v1/events", authMiddleware(http.HandlerFunc(eventHandler.CreateEvent)))
	mux.Handle("PUT /api/v1/events/{id}", authMiddleware(http.HandlerFunc(eventHandler.UpdateEvent)))
	mux.Handle("DELETE /api/v1/events/{id}", authMiddleware(http.HandlerFunc(eventHandler.DeleteEvent)))

	mux.HandleFunc("GET /api/v1/events", eventHandler.ListEvents)
	mux.HandleFunc("GET /api/v1/events/{id}", eventHandler.GetEvent)
	mux.HandleFunc("GET /api/v1/events/{id}/seats", eventHandler.GetAvailableSeats)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	go func() {
		log.Info("HTTP server is running", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()

	log.Info("API Gateway is running", "port", cfg.Server.Port)
	log.Info("Press Ctrl+C to stop")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down...")
	server.Close()
	log.Info("API Gateway stopped")
}

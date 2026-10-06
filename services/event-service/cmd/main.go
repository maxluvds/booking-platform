package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/maxluvds/booking-platform/api/proto"
	"github.com/maxluvds/booking-platform/pkg/config"
	"github.com/maxluvds/booking-platform/pkg/database"
	"github.com/maxluvds/booking-platform/pkg/logger"
	"github.com/maxluvds/booking-platform/services/event-service/internal/delivery"
	"github.com/maxluvds/booking-platform/services/event-service/internal/repository"
	"github.com/maxluvds/booking-platform/services/event-service/internal/usecase"
)

func main() {
	cfg, err := config.Load("services/event-service/config/config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	log := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)

	log.Info("Event Service starting...")

	log.Info("Running database migrations...")
	if err := database.RunMigrations(
		cfg.Database.URL,
		"./migrations/event-service",
	); err != nil {
		log.Error("Failed to run migrations", "error", err)
		return
	}
	log.Info("Migrations completed successfully")

	log.Info("Connecting to database...")
	db, err := database.NewPool(
		cfg.Database.GetDSN(),
		cfg.Database.MaxConnections,
	)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()
	log.Info("Database connected successfully")

	eventRepo := repository.NewEventRepository(db.Conn())
	eventUsecase := usecase.NewEventUseCase(eventRepo)
	grpcHandler := delivery.NewEventGRPCHandler(eventUsecase)

	log.Info("Server settings",
		"port", cfg.Server.Port,
		"timeout", cfg.Server.Timeout,
	)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Error("Failed to listen", "error", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEventServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	go func() {
		log.Info("gRPC server is running", "port", cfg.Server.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("Failed to serve gRPC", "error", err)
		}
	}()

	log.Info("Event Service is running", "port", cfg.Server.Port)
	log.Info("Press Ctrl+C to stop")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Info("Event Service stopped")
}


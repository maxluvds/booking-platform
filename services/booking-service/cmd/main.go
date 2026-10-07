package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/maxluvds/booking-platform/api/proto"
	"github.com/maxluvds/booking-platform/pkg/config"
	"github.com/maxluvds/booking-platform/pkg/database"
	"github.com/maxluvds/booking-platform/pkg/kafka"
	"github.com/maxluvds/booking-platform/pkg/logger"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/delivery"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/repository"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/usecase"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/worker"
)

func main() {
	cfg, err := config.Load("services/booking-service/config/config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	log := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)

	log.Info("Booking Service starting...")

	log.Info("Running database migrations...")
	if err := database.RunMigrations(
		cfg.Database.URL,
		"./migrations/booking-service",
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

	log.Info("Connecting to Kafka...")
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, kafka.TopicBookingCreated)
	defer kafkaProducer.Close()
	log.Info("Kafka producer created successfully")

	bookingRepo := repository.NewBookingRepository(db.Conn())
	outboxRepo := repository.NewOutboxRepository(db.Conn())
	bookingUsecase := usecase.NewBookingUseCase(bookingRepo, outboxRepo)
	grpcHandler := delivery.NewBookingGRPCHandler(bookingUsecase)

	outboxWorker := worker.NewOutboxWorker(
		outboxRepo,
		kafkaProducer,
		log,
		cfg.Outbox.Interval,
		cfg.Outbox.BatchSize,
	)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	go outboxWorker.Start(workerCtx)

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
	pb.RegisterBookingServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	go func() {
		log.Info("gRPC server is running", "port", cfg.Server.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("Failed to serve gRPC", "error", err)
		}
	}()

	log.Info("Booking Service is running", "port", cfg.Server.Port)
	log.Info("Press Ctrl+C to stop")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down...")
	workerCancel()
	grpcServer.GracefulStop()
	log.Info("Booking Service stopped")
}

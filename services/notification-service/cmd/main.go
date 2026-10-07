package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/maxluvds/booking-platform/pkg/config"
	"github.com/maxluvds/booking-platform/pkg/database"
	"github.com/maxluvds/booking-platform/pkg/kafka"
	"github.com/maxluvds/booking-platform/pkg/logger"
	"github.com/maxluvds/booking-platform/services/notification-service/internal/consumer"
	"github.com/maxluvds/booking-platform/services/notification-service/internal/repository"
	"github.com/maxluvds/booking-platform/services/notification-service/internal/usecase"
)

func main() {
	cfg, err := config.Load("services/notification-service/config/config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	log := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)

	log.Info("Notification Service starting...")

	log.Info("Running database migrations...")
	if err := database.RunMigrations(
		cfg.Database.URL,
		"./migrations/notification-service",
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
	kafkaConsumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		kafka.TopicEventCreated,
		cfg.Kafka.GroupID,
	)
	defer kafkaConsumer.Close()
	log.Info("Kafka consumer created successfully")

	notificationRepo := repository.NewNotificationRepository(db.Conn())
	notificationUsecase := usecase.NewNotificationUseCase(notificationRepo)
	kafkaHandler := consumer.NewKafkaConsumer(kafkaConsumer, notificationUsecase, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaHandler.Start(ctx)

	log.Info("Notification Service is running")
	log.Info("Press Ctrl+C to stop")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down...")
	cancel()
	log.Info("Notification Service stopped")
}

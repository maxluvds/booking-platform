package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/maxluvds/booking-platform/pkg/kafka"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/repository"
)

type OutboxWorker struct {
	outboxRepo    *repository.OutboxRepository
	kafkaProducer *kafka.Producer
	log           *slog.Logger
	interval      time.Duration
	batchSize     int
}

func NewOutboxWorker(
	outboxRepo *repository.OutboxRepository,
	kafkaProducer *kafka.Producer,
	log *slog.Logger,
	interval time.Duration,
	batchSize int,
) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo:    outboxRepo,
		kafkaProducer: kafkaProducer,
		log:           log,
		interval:      interval,
		batchSize:     batchSize,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	w.log.Info("Outbox worker started",
		"interval", w.interval,
		"batch_size", w.batchSize,
	)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info("Outbox worker stopping...")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context) {
	events, err := w.outboxRepo.GetPending(ctx, w.batchSize)
	if err != nil {
		w.log.Error("Failed to get pending events", "error", err)
		return
	}

	if len(events) == 0 {
		return
	}

	w.log.Info("Processing outbox events", "count", len(events))

	for _, event := range events {
		if err := w.kafkaProducer.SendRaw(ctx, event.EventType, event.Payload); err != nil {
			w.log.Error("Failed to send event to Kafka",
				"event_id", event.ID,
				"event_type", event.EventType,
				"error", err,
			)

			if err := w.outboxRepo.MarkAsFailed(ctx, event.ID, err.Error()); err != nil {
				w.log.Error("Failed to mark event as failed",
					"event_id", event.ID,
					"error", err,
				)
			}
			continue
		}

		if err := w.outboxRepo.MarkAsSent(ctx, event.ID); err != nil {
			w.log.Error("Failed to mark event as sent",
				"event_id", event.ID,
				"error", err,
			)
			continue
		}

		w.log.Info("Event sent successfully",
			"event_id", event.ID,
			"event_type", event.EventType,
		)
	}
}

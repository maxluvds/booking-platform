package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/maxluvds/booking-platform/pkg/kafka"
	"github.com/maxluvds/booking-platform/services/event-service/internal/domain"
	"github.com/maxluvds/booking-platform/services/event-service/internal/repository"
)

type EventUseCase struct {
	repo          *repository.EventRepository
	kafkaProducer *kafka.Producer
}

func NewEventUseCase(repo *repository.EventRepository, kafkaProducer *kafka.Producer) *EventUseCase {
	return &EventUseCase{
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

func (uc *EventUseCase) Create(ctx context.Context, req *domain.CreateEventRequest) (*domain.Event, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.City == "" {
		return nil, fmt.Errorf("city is required")
	}
	if req.TotalSeats <= 0 {
		return nil, fmt.Errorf("total seats must be greater than 0")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	event, err := uc.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	message := kafka.EventCreatedMessage{
		EventID:   event.ID,
		Title:     event.Title,
		City:      event.City,
		Category:  event.Category,
		Date:      event.Date,
		Price:     event.Price,
		Timestamp: time.Now(),
	}

	if err := uc.kafkaProducer.Send(ctx, fmt.Sprintf("%d", event.ID), message); err != nil {
		fmt.Printf("Failed to send Kafka message: %v\n", err)
	}

	return event, nil
}

func (uc *EventUseCase) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid event id")
	}

	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	return event, nil
}

func (uc *EventUseCase) List(ctx context.Context, req *domain.ListEventsRequest) ([]*domain.Event, int, error) {
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	return uc.repo.List(ctx, req)
}

func (uc *EventUseCase) Update(ctx context.Context, id int64, req *domain.UpdateEventRequest) (*domain.Event, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid event id")
	}

	event, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	return event, nil
}

func (uc *EventUseCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid event id")
	}

	return uc.repo.Delete(ctx, id)
}
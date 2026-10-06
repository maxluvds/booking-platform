package usecase

import (
	"context"
	"fmt"

	"github.com/maxluvds/booking-platform/services/event-service/internal/domain"
	"github.com/maxluvds/booking-platform/services/event-service/internal/repository"
)

type EventUseCase struct {
	repo *repository.EventRepository
}

func NewEventUseCase(repo *repository.EventRepository) *EventUseCase {
	return &EventUseCase{repo: repo}
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

	return uc.repo.Create(ctx, req)
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

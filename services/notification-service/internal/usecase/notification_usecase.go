package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/maxluvds/booking-platform/pkg/kafka"
	"github.com/maxluvds/booking-platform/services/notification-service/internal/domain"
	"github.com/maxluvds/booking-platform/services/notification-service/internal/repository"
)

type NotificationUseCase struct {
	repo *repository.NotificationRepository
}

func NewNotificationUseCase(repo *repository.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

func (uc *NotificationUseCase) HandleEventCreated(ctx context.Context, msg kafka.EventCreatedMessage) error {
	title := fmt.Sprintf("New event: %s", msg.Title)
	message := fmt.Sprintf(
		"Event '%s' in %s on %s. Price: %.2f",
		msg.Title,
		msg.City,
		msg.Date.Format(time.RFC3339),
		msg.Price,
	)

	notification := &domain.Notification{
		EventID: msg.EventID,
		Type:    "event.created",
		Title:   title,
		Message: message,
		Status:  "sent",
	}

	_, err := uc.repo.Create(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	return nil
}

func (uc *NotificationUseCase) List(ctx context.Context, limit, offset int) ([]*domain.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return uc.repo.List(ctx, limit, offset)
}

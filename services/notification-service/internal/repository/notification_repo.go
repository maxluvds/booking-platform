package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxluvds/booking-platform/services/notification-service/internal/domain"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	query := `
		INSERT INTO notifications (event_id, type, title, message, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, event_id, type, title, message, status, created_at
	`

	var notification domain.Notification
	err := r.db.QueryRow(
		ctx, query,
		n.EventID,
		n.Type,
		n.Title,
		n.Message,
		n.Status,
		time.Now(),
	).Scan(
		&notification.ID,
		&notification.EventID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.Status,
		&notification.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

func (r *NotificationRepository) List(ctx context.Context, limit, offset int) ([]*domain.Notification, error) {
	query := `
		SELECT id, event_id, type, title, message, status, created_at
		FROM notifications
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(
			&n.ID,
			&n.EventID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.Status,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, &n)
	}

	return notifications, nil
}

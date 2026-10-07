package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxluvds/booking-platform/services/booking-service/internal/domain"
)

type OutboxRepository struct {
	db *pgxpool.Pool
}

func NewOutboxRepository(db *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) GetPending(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	query := `
		SELECT id, aggregate_id, event_type, payload, status, retry_count, error_message, created_at, sent_at
		FROM outbox
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending events: %w", err)
	}
	defer rows.Close()

	var events []*domain.OutboxEvent
	for rows.Next() {
		var e domain.OutboxEvent
		err := rows.Scan(
			&e.ID,
			&e.AggregateID,
			&e.EventType,
			&e.Payload,
			&e.Status,
			&e.RetryCount,
			&e.ErrorMessage,
			&e.CreatedAt,
			&e.SentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &e)
	}

	return events, nil
}

func (r *OutboxRepository) MarkAsSent(ctx context.Context, id int64) error {
	query := `
		UPDATE outbox
		SET status = 'sent', sent_at = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark as sent: %w", err)
	}

	return nil
}

func (r *OutboxRepository) MarkAsFailed(ctx context.Context, id int64, errMsg string) error {
	query := `
		UPDATE outbox
		SET status = 'failed',
		    retry_count = retry_count + 1,
		    error_message = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to mark as failed: %w", err)
	}

	return nil
}

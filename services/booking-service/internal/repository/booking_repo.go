package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxluvds/booking-platform/services/booking-service/internal/domain"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(
	ctx context.Context,
	booking *domain.Booking,
	outboxEvent *domain.OutboxEvent,
) (*domain.Booking, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	bookingQuery := `
		INSERT INTO bookings (user_id, event_id, seats, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, event_id, seats, total_price, status, created_at, updated_at
	`

	var created domain.Booking
	err = tx.QueryRow(
		ctx, bookingQuery,
		booking.UserID,
		booking.EventID,
		booking.Seats,
		booking.TotalPrice,
		booking.Status,
		time.Now(),
		time.Now(),
	).Scan(
		&created.ID,
		&created.UserID,
		&created.EventID,
		&created.Seats,
		&created.TotalPrice,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	outboxQuery := `
		INSERT INTO outbox (aggregate_id, event_type, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = tx.Exec(
		ctx, outboxQuery,
		created.ID,
		outboxEvent.EventType,
		outboxEvent.Payload,
		"pending",
		time.Now(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create outbox event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &created, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, event_id, seats, total_price, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	var b domain.Booking
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.EventID,
		&b.Seats,
		&b.TotalPrice,
		&b.Status,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &b, nil
}

func (r *BookingRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.Booking, error) {
	query := `
		SELECT id, user_id, event_id, seats, total_price, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.EventID,
			&b.Seats,
			&b.TotalPrice,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &b)
	}

	return bookings, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id int64, status domain.BookingStatus) error {
	query := `
		UPDATE bookings
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

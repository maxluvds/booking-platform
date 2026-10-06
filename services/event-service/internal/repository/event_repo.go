package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxluvds/booking-platform/services/event-service/internal/domain"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, req *domain.CreateEventRequest) (*domain.Event, error) {
	query := `
		INSERT INTO events (
			title, description, city, address, category, date, 
			total_seats, available_seats, price, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, title, description, city, address, category, date, 
			total_seats, available_seats, price, status, created_at, updated_at
	`

	var event domain.Event
	err := r.db.QueryRow(
		ctx, query,
		req.Title,
		req.Description,
		req.City,
		req.Address,
		req.Category,
		req.Date,
		req.TotalSeats,
		req.TotalSeats,
		req.Price,
		domain.StatusActive,
		time.Now(),
		time.Now(),
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.City,
		&event.Address,
		&event.Category,
		&event.Date,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.Price,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	query := `
		SELECT id, title, description, city, address, category, date, 
			total_seats, available_seats, price, status, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event domain.Event
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.City,
		&event.Address,
		&event.Category,
		&event.Date,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.Price,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) List(ctx context.Context, req *domain.ListEventsRequest) ([]*domain.Event, int, error) {
	conditions := []string{}
	args := []interface{}{}
	argCounter := 1

	if req.City != "" {
		conditions = append(conditions, fmt.Sprintf("city = $%d", argCounter))
		args = append(args, req.City)
		argCounter++
	}

	if req.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argCounter))
		args = append(args, req.Category)
		argCounter++
	}

	if req.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, req.Status)
		argCounter++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM events %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, title, description, city, address, category, date, 
			total_seats, available_seats, price, status, created_at, updated_at
		FROM events
		%s
		ORDER BY date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCounter, argCounter+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.City,
			&event.Address,
			&event.Category,
			&event.Date,
			&event.TotalSeats,
			&event.AvailableSeats,
			&event.Price,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, total, nil
}

func (r *EventRepository) Update(ctx context.Context, id int64, req *domain.UpdateEventRequest) (*domain.Event, error) {
	updates := []string{}
	args := []interface{}{}
	argCounter := 1

	if req.Title != nil {
		updates = append(updates, fmt.Sprintf("title = $%d", argCounter))
		args = append(args, *req.Title)
		argCounter++
	}

	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argCounter))
		args = append(args, *req.Description)
		argCounter++
	}

	if req.City != nil {
		updates = append(updates, fmt.Sprintf("city = $%d", argCounter))
		args = append(args, *req.City)
		argCounter++
	}

	if req.Address != nil {
		updates = append(updates, fmt.Sprintf("address = $%d", argCounter))
		args = append(args, *req.Address)
		argCounter++
	}

	if req.Category != nil {
		updates = append(updates, fmt.Sprintf("category = $%d", argCounter))
		args = append(args, *req.Category)
		argCounter++
	}

	if req.Date != nil {
		updates = append(updates, fmt.Sprintf("date = $%d", argCounter))
		args = append(args, *req.Date)
		argCounter++
	}

	if req.TotalSeats != nil {
		updates = append(updates, fmt.Sprintf("total_seats = $%d", argCounter))
		args = append(args, *req.TotalSeats)
		argCounter++
	}

	if req.Price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argCounter))
		args = append(args, *req.Price)
		argCounter++
	}

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, *req.Status)
		argCounter++
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCounter))
	args = append(args, time.Now())
	argCounter++

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE events
		SET %s
		WHERE id = $%d
		RETURNING id, title, description, city, address, category, date, 
			total_seats, available_seats, price, status, created_at, updated_at
	`, strings.Join(updates, ", "), argCounter)

	var event domain.Event
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.City,
		&event.Address,
		&event.Category,
		&event.Date,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.Price,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *EventRepository) UpdateAvailableSeats(ctx context.Context, eventID int64, seats int) error {
	query := `
		UPDATE events 
		SET available_seats = available_seats - $1, updated_at = $2
		WHERE id = $3 AND available_seats >= $1
	`

	result, err := r.db.Exec(ctx, query, seats, time.Now(), eventID)
	if err != nil {
		return fmt.Errorf("failed to update available seats: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("not enough seats available")
	}

	return nil
}

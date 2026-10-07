package domain

import "time"

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID         int64         `json:"id"`
	UserID     int64         `json:"user_id"`
	EventID    int64         `json:"event_id"`
	Seats      int           `json:"seats"`
	TotalPrice float64       `json:"total_price"`
	Status     BookingStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type CreateBookingRequest struct {
	UserID  int64 `json:"user_id"`
	EventID int64 `json:"event_id"`
	Seats   int   `json:"seats"`
}

type CancelBookingRequest struct {
	BookingID int64 `json:"booking_id"`
	UserID    int64 `json:"user_id"`
}

type OutboxEvent struct {
	ID           int64      `json:"id"`
	AggregateID  int64      `json:"aggregate_id"`
	EventType    string     `json:"event_type"`
	Payload      []byte     `json:"payload"`
	Status       string     `json:"status"`
	RetryCount   int        `json:"retry_count"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
}

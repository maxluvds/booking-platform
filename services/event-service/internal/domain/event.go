package domain

import (
	"time"
)

type EventStatus string

const (
	StatusActive   EventStatus = "active"
	StatusCanceled EventStatus = "canceled"
	StatusFinished EventStatus = "finished"
)

type Event struct {
	ID             int64       `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	City           string      `json:"city"`
	Address        string      `json:"address"`
	Category       string      `json:"category"`
	Date           time.Time   `json:"date"`
	TotalSeats     int         `json:"total_seats"`
	AvailableSeats int         `json:"available_seats"`
	Price          float64     `json:"price"`
	Status         EventStatus `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	City        string    `json:"city" validate:"required,min=2,max=100"`
	Address     string    `json:"address" validate:"required,min=2,max=255"`
	Category    string    `json:"category" validate:"required,min=2,max=50"`
	Date        time.Time `json:"date" validate:"required"`
	TotalSeats  int       `json:"total_seats" validate:"required,min=1"`
	Price       float64   `json:"price" validate:"required,min=0"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	City        *string    `json:"city,omitempty"`
	Address     *string    `json:"address,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	TotalSeats  *int       `json:"total_seats,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	Status      *string    `json:"status,omitempty"`
}

type ListEventsRequest struct {
	City     string `json:"city,omitempty"`
	Category string `json:"category,omitempty"`
	Status   string `json:"status,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

package kafka

import "time"

const (
	TopicEventCreated   = "event.created"
	TopicEventUpdated   = "event.updated"
	TopicEventDeleted   = "event.deleted"
	TopicBookingCreated = "booking.created"
)

type EventCreatedMessage struct {
	EventID   int64     `json:"event_id"`
	Title     string    `json:"title"`
	City      string    `json:"city"`
	Category  string    `json:"category"`
	Date      time.Time `json:"date"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type EventUpdatedMessage struct {
	EventID   int64     `json:"event_id"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
}

type EventDeletedMessage struct {
	EventID   int64     `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
}

package domain

import "time"

type Notification struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

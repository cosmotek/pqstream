package pqstream

import "time"

type Event struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}

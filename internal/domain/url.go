package domain

import "time"

type URL struct {
	Original  string    `json:"original"`
	Shortened string    `json:"shortened"`
	CreatedAt time.Time `json:"created_at"`
}

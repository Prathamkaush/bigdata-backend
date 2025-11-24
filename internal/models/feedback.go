package models

import "time"

type Feedback struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

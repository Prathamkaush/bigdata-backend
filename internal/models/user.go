package models

import "time"

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	APIKeyHash    string    `json:"-"`                 // stored hash
	APIKey        string    `json:"api_key,omitempty"` // only shown when created/regenerated
	Credits       int       `json:"credits"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	TotalRequests int       `json:"total_requests,omitempty"`
	TodayRequests int       `json:"today_requests,omitempty"`
	CreditsUsed   int       `json:"credits_used,omitempty"`
}

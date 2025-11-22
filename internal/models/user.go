package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	ApiKey    string    `json:"api_key"`
	Credits   int64     `json:"credits"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}

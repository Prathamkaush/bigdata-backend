package models

type ApiLog struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Endpoint   string `json:"endpoint"`
	DurationMs int64  `json:"duration_ms"`
	Records    int64  `json:"records"`
	CreatedAt  string `json:"created_at"`
}

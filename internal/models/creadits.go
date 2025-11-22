package models

type CreditLog struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Used      int64  `json:"used"`
	Endpoint  string `json:"endpoint"`
	Params    string `json:"params"`
	CreatedAt string `json:"created_at"`
}

package domain

import "time"

type Story struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	UserName   string    `json:"user_name" db:"user_name"`
	UserAvatar *string   `json:"user_avatar" db:"user_avatar"`
	FileURL    string    `json:"file_url" db:"file_url"`
	Type       string    `json:"type" db:"type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	Viewed     bool      `json:"viewed" db:"viewed"`
}

package domain

import "time"

type Chat struct {
	ID          int64     `json:"id" db:"id"`
	Type        string    `json:"type" db:"type" binding:"required,oneof=private group channel"`
	Name        string    `json:"name" db:"name" binding:"required,min=2,max=255"`
	Description string    `json:"description,omitempty" db:"name" binding:"max=500"`
	AvatarPath  *string   `json:"avatar_path,omitempty" db:"avatar_path"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
type ChatResponse struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	AvatarPath  *string   `json:"avatar_path,omitempty"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type CreateChatRequest struct {
	Type string `json:"type" db:"type" binding:"required,oneof=private group channel"`
	Name string `json:"name" db:"name" binding:"required,min=2,max=255"`
}

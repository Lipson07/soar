package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ChatID    uuid.UUID  `json:"chat_id" db:"chat_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Text      string     `json:"text" db:"text"`
	ReplyTo   *uuid.UUID `json:"reply_to,omitempty" db:"reply_to"`
	IsEdited  bool       `json:"is_edited" db:"is_edited"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type SendMessageRequest struct {
	Text    string     `json:"text" validate:"required"`
	ReplyTo *uuid.UUID `json:"reply_to,omitempty"`
}

type MessageInfo struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateMessageRequest struct {
	Text string `json:"text" validate:"required"`
}
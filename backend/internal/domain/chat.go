package domain

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Name          *string    `json:"name,omitempty" db:"name"`
	Type          ChatType   `json:"type" db:"type"`
	CreatorID     uuid.UUID  `json:"creator_id" db:"creator_id"`
	AvatarURL     *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty" db:"last_message_at"`
}

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
)

type CreatePrivateChatRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type CreateGroupChatRequest struct {
	Name      string      `json:"name" validate:"required"`
	UserIDs   []uuid.UUID `json:"user_ids" validate:"required,min=1"`
	AvatarURL *string     `json:"avatar_url,omitempty"`
}

type UpdateChatRequest struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type ChatResponse struct {
	ID            uuid.UUID    `json:"id"`
	Name          *string      `json:"name,omitempty"`
	Type          ChatType     `json:"type"`
	CreatorID     uuid.UUID    `json:"creator_id"`
	AvatarURL     *string      `json:"avatar_url,omitempty"`
	LastMessage   *MessageInfo `json:"last_message,omitempty"`
	UnreadCount   int          `json:"unread_count"`
	LastMessageAt *time.Time   `json:"last_message_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
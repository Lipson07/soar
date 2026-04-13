package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	ChatID    uuid.UUID   `json:"chat_id" db:"chat_id"`
	UserID    uuid.UUID   `json:"user_id" db:"user_id"`
	Type      MessageType `json:"type" db:"type"`
	Text      string      `json:"text,omitempty" db:"text"`
	FileURL   string      `json:"file_url,omitempty" db:"file_url"`
	FileName  string      `json:"file_name,omitempty" db:"file_name"`
	FileSize  int64       `json:"file_size,omitempty" db:"file_size"`
	MimeType  string      `json:"mime_type,omitempty" db:"mime_type"`
	ReplyTo   *uuid.UUID  `json:"reply_to,omitempty" db:"reply_to"`
	IsEdited  bool        `json:"is_edited" db:"is_edited"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`
}

type SendMessageRequest struct {
	Type     MessageType `json:"type" validate:"required,oneof=text image file"`
	Text     string      `json:"text,omitempty"`
	FileURL  string      `json:"file_url,omitempty"`
	FileName string      `json:"file_name,omitempty"`
	FileSize int64       `json:"file_size,omitempty"`
	MimeType string      `json:"mime_type,omitempty"`
	ReplyTo  *uuid.UUID  `json:"reply_to,omitempty"`
}

type MessageInfo struct {
	ID        uuid.UUID   `json:"id"`
	Type      MessageType `json:"type"`
	Text      string      `json:"text,omitempty"`
	FileURL   string      `json:"file_url,omitempty"`
	FileName  string      `json:"file_name,omitempty"`
	UserID    uuid.UUID   `json:"user_id"`
	CreatedAt time.Time   `json:"created_at"`
}

type UpdateMessageRequest struct {
	Text string `json:"text" validate:"required"`
}

type UploadFileResponse struct {
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}
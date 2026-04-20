package domain

import (
	"time"

	"github.com/google/uuid"
)

type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeFile  FileType = "file"
)

type File struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	URL       string     `json:"url" db:"url"`
	Size      int64      `json:"size" db:"size"`
	Type      FileType   `json:"type" db:"type"`
	MimeType  string     `json:"mime_type" db:"mime_type"`
	ChatID    *uuid.UUID `json:"chat_id,omitempty" db:"chat_id"`
	MessageID *uuid.UUID `json:"message_id,omitempty" db:"message_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type FileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	Type      string `json:"type"`
	MimeType  string `json:"mime_type"`
	CreatedAt string `json:"created_at"`
}

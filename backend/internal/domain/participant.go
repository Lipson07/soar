package domain

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ChatID     uuid.UUID `json:"chat_id" db:"chat_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Role       string    `json:"role" db:"role"` // admin, member
	JoinedAt   time.Time `json:"joined_at" db:"joined_at"`
	LastReadAt time.Time `json:"last_read_at" db:"last_read_at"`
}

type ParticipantRole string

const (
	RoleAdmin  ParticipantRole = "admin"
	RoleMember ParticipantRole = "member"
)

type AddParticipantsRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" validate:"required,min=1"`
}

type UpdateRoleRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Role   string    `json:"role" validate:"required,oneof=admin member"`
}

type ParticipantResponse struct {
	ID         uuid.UUID `json:"id"`
	ChatID     uuid.UUID `json:"chat_id"`
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joined_at"`
	LastReadAt time.Time `json:"last_read_at"`
}
package domain

import "time"

type ChatMember struct {
	ID         int64     `json:"id" db:"id"`
	ChatID     int64     `json:"chat_id" db:"chat_id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Role       string    `json:"role" db:"role" binding:"required,oneof=owner admin member"`
	JoinedAt   time.Time `json:"joined_at" db:"joined_at"`
	LastReadAt time.Time `json:"last_read_at" db:"last_read_at"`
}

type AddMemberRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"omitempty,oneof=admin member"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"`
}

type RemoveMemberRequest struct {
	ChatID int64 `json:"chat_id" binding:"required"`
	UserID int64 `json:"user_id" binding:"required"`
}

type LeaveChatRequest struct {
	ChatID int64 `json:"chat_id" binding:"required"`
}

type KickMemberRequest struct {
	ChatID int64 `json:"chat_id" binding:"required"`
	UserID int64 `json:"user_id" binding:"required"`
}

type DeleteMemberResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

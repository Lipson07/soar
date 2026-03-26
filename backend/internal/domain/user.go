package domain

import (
	"time"
)

type User struct {
	ID              int64      `json:"id" db:"id"`
	Name            string     `json:"name" db:"name" binding:"required,min=2,max=255"`
	Email           string     `json:"email" db:"email" binding:"required,email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	Password        string     `json:"-" db:"password"`
	Role            string     `json:"role" db:"role" binding:"required"`
	Avatar          bool       `json:"avatar" db:"avatar"`
	AvatarPath      *string    `json:"avatar_path,omitempty" db:"avatar_path"`
	LastSeenAt      *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
	IsOnline        bool       `json:"is_online" db:"is_online"`
	RememberToken   *string    `json:"-" db:"remember_token"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user admin moderator"`
}

type UpdateUserRequest struct {
	Name       *string `json:"name,omitempty" binding:"omitempty,min=2"`
	Email      *string `json:"email,omitempty" binding:"omitempty,email"`
	Role       *string `json:"role,omitempty" binding:"omitempty,oneof=user admin moderator"`
	Avatar     *bool   `json:"avatar,omitempty"`
	AvatarPath *string `json:"avatar_path,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	Role            string     `json:"role"`
	Avatar          bool       `json:"avatar"`
	AvatarPath      *string    `json:"avatar_path,omitempty"`
	LastSeenAt      *time.Time `json:"last_seen_at,omitempty"`
	IsOnline        bool       `json:"is_online"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
type PasswordResetToken struct {
	Email     string    `json:"email" db:"email"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	Payload      string    `json:"payload" db:"payload"`
	LastActivity int64     `json:"last_activity" db:"last_activity"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
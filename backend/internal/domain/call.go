package domain

import (
	"time"

	"github.com/google/uuid"
)

type CallType string

const (
	CallTypeAudio CallType = "audio"
	CallTypeVideo CallType = "video"
)

type CallStatus string

const (
	CallStatusPending  CallStatus = "pending"
	CallStatusActive   CallStatus = "active"
	CallStatusRejected CallStatus = "rejected"
	CallStatusMissed   CallStatus = "missed"
	CallStatusEnded    CallStatus = "ended"
)

type Call struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ChatID    uuid.UUID  `json:"chat_id" db:"chat_id"`
	CallerID  uuid.UUID  `json:"caller_id" db:"caller_id"`
	CalleeID  uuid.UUID  `json:"callee_id" db:"callee_id"`
	Type      CallType   `json:"type" db:"type"`
	Status    CallStatus `json:"status" db:"status"`
	RoomID    string     `json:"room_id" db:"room_id"`
	StartedAt *time.Time `json:"started_at,omitempty" db:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type CallSignal struct {
	Type      string        `json:"type"`
	CallID    uuid.UUID     `json:"call_id"`
	RoomID    string        `json:"room_id"`
	ChatID    string        `json:"chat_id,omitempty"`
	CalleeID  string        `json:"callee_id,omitempty"`
	CallType  string        `json:"call_type,omitempty"`
	FromID    uuid.UUID     `json:"from_id"`
	ToID      uuid.UUID     `json:"to_id"`
	SDP       string        `json:"sdp,omitempty"`
	Candidate interface{}   `json:"candidate,omitempty"`
	Call      *CallResponse `json:"call,omitempty"`
}

type CallResponse struct {
	ID       uuid.UUID  `json:"id"`
	ChatID   uuid.UUID  `json:"chat_id"`
	CallerID uuid.UUID  `json:"caller_id"`
	CalleeID uuid.UUID  `json:"callee_id"`
	Type     CallType   `json:"type"`
	Status   CallStatus `json:"status"`
	RoomID   string     `json:"room_id"`
	Caller   *UserInfo  `json:"caller,omitempty"`
}

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

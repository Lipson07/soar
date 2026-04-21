package repository

import (
	"backend/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Search(ctx context.Context, query string) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type ChatRepository interface {
	Create(ctx context.Context, chat *domain.Chat) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
	GetPrivateChatByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*domain.Chat, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error)
	GetAll(ctx context.Context) ([]*domain.Chat, error)
	Update(ctx context.Context, chat *domain.Chat) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastMessage(ctx context.Context, chatID uuid.UUID, lastMessageAt time.Time) error
}

type ParticipantRepository interface {
	Add(ctx context.Context, participant *domain.Participant) error
	Remove(ctx context.Context, chatID, userID uuid.UUID) error
	GetByChatID(ctx context.Context, chatID uuid.UUID) ([]*domain.Participant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Participant, error)
	IsParticipant(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	IsAdmin(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	GetByChatAndUser(ctx context.Context, chatID, userID uuid.UUID) (*domain.Participant, error)
	Update(ctx context.Context, participant *domain.Participant) error
	UpdateLastRead(ctx context.Context, chatID, userID uuid.UUID, lastReadAt time.Time) error
	GetUnreadCount(ctx context.Context, chatID, userID uuid.UUID, lastReadAt time.Time) (int, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*domain.Message, error)
	Update(ctx context.Context, message *domain.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetLastMessage(ctx context.Context, chatID uuid.UUID) (*domain.MessageInfo, error)
}
type SecurityRepository interface {
	GetUserSecuritySettings(ctx context.Context, userID uuid.UUID) (*domain.SecuritySettings, error)
	CreateDefaultSettings(ctx context.Context, userID uuid.UUID) (*domain.SecuritySettings, error)
	UpdateSecuritySettings(ctx context.Context, settings *domain.SecuritySettings) error
	EnableTwoFactor(ctx context.Context, userID uuid.UUID, secret string) error
	DisableTwoFactor(ctx context.Context, userID uuid.UUID) error
	CreateAuditLog(ctx context.Context, log *domain.SecurityAuditLog) error
	GetUserAuditLogs(ctx context.Context, userID uuid.UUID, limit int) ([]domain.SecurityAuditLog, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.UserSession) error
	GetUserActiveSessions(ctx context.Context, userID uuid.UUID) ([]domain.UserSession, error)
	TerminateSession(ctx context.Context, sessionID int64, userID uuid.UUID) error
	TerminateAllOtherSessions(ctx context.Context, userID uuid.UUID, currentSessionToken string) error
	UpdateSessionActivity(ctx context.Context, sessionToken string) error
	GetSessionByToken(ctx context.Context, token string) (*domain.UserSession, error)
}
type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error)
	GetAll(ctx context.Context) ([]*domain.File, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.File, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID) ([]*domain.File, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
type CallRepository interface {
	Create(ctx context.Context, call *domain.Call) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Call, error)
	GetByRoomID(ctx context.Context, roomID string) (*domain.Call, error)
	GetActiveCall(ctx context.Context, chatID uuid.UUID) (*domain.Call, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CallStatus) error
	UpdateEnded(ctx context.Context, id uuid.UUID, endedAt time.Time) error
	GetUserCalls(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Call, error)
}

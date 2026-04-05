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
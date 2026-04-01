package repository

import (
	"context"
	"myapp/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, query string, limit, offset int) ([]domain.User, error)
}
type ChatRepository interface {
	Create(ctx context.Context, user *domain.Chat) error
	GetByID(ctx context.Context, id int64) (*domain.Chat, error)
	GetByName(ctx context.Context, email string) (*domain.Chat, error)
	GetAll(ctx context.Context) ([]domain.Chat, error)
	Update(ctx context.Context, user *domain.Chat) error
	Delete(ctx context.Context, id int64) error
}
type ChatMemberRepository interface {
	Create(ctx context.Context, member *domain.ChatMember) error
	GetByID(ctx context.Context, id int64) (*domain.ChatMember, error)
	GetByChatAndUser(ctx context.Context, chatID, userID int64) (*domain.ChatMember, error)
	GetByChatID(ctx context.Context, chatID int64) ([]domain.ChatMember, error)
	GetByUserID(ctx context.Context, userID int64) ([]domain.ChatMember, error)
	Update(ctx context.Context, member *domain.ChatMember) error
	UpdateRole(ctx context.Context, chatID, userID int64, role string) error
	UpdateLastReadAt(ctx context.Context, chatID, userID int64) error
	Delete(ctx context.Context, id int64) error
	DeleteByChatAndUser(ctx context.Context, chatID, userID int64) error
	DeleteAllByChatID(ctx context.Context, chatID int64) error
	Exists(ctx context.Context, chatID, userID int64) (bool, error)
	CountMembers(ctx context.Context, chatID int64) (int, error)
	GetUserRole(ctx context.Context, chatID, userID int64) (string, error)
	IsMember(ctx context.Context, chatID, userID int64) (bool, error)
	IsOwner(ctx context.Context, chatID, userID int64) (bool, error)
	IsAdmin(ctx context.Context, chatID, userID int64) (bool, error)
}

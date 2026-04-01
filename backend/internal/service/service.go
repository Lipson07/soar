package service

import (
	"context"
	"myapp/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]domain.User, error)
}
type ChatService interface {
	Create(ctx context.Context, req *domain.CreateChatRequest) (*domain.Chat, error)
	GetByID(ctx context.Context, id int64) (*domain.Chat, error)
	GetByName(ctx context.Context, name string) (*domain.Chat, error)
	GetAll(ctx context.Context) ([]domain.Chat, error)
	Update(ctx context.Context, id int64, req *domain.UpdateChatRequest) (*domain.Chat, error)
	Delete(ctx context.Context, id int64) error
}

type ChatMemberService interface {
	AddMember(ctx context.Context, chatID int64, req *domain.AddMemberRequest, currentUserID int64) (*domain.ChatMember, error)
	AddMembers(ctx context.Context, chatID int64, userIDs []int64, currentUserID int64) error
	GetMember(ctx context.Context, chatID, userID int64) (*domain.ChatMember, error)
	GetChatMembers(ctx context.Context, chatID int64) ([]domain.ChatMember, error)
	GetUserChats(ctx context.Context, userID int64) ([]domain.ChatMember, error)
	GetMemberCount(ctx context.Context, chatID int64) (int, error)
	UpdateMemberRole(ctx context.Context, chatID, userID int64, req *domain.UpdateMemberRoleRequest, currentUserID int64) error
	UpdateLastRead(ctx context.Context, chatID, userID int64) error
	RemoveMember(ctx context.Context, chatID, userID int64, currentUserID int64) error
	LeaveChat(ctx context.Context, chatID, userID int64) error
	KickMember(ctx context.Context, chatID, userID int64, currentUserID int64) error
	IsMember(ctx context.Context, chatID, userID int64) (bool, error)
	IsAdmin(ctx context.Context, chatID, userID int64) (bool, error)
	IsOwner(ctx context.Context, chatID, userID int64) (bool, error)
	GetUserRole(ctx context.Context, chatID, userID int64) (string, error)
}

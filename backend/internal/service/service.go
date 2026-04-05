package service

import (
	"backend/internal/domain"
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string) ([]*domain.User, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type ChatService interface {
	CreatePrivateChat(ctx context.Context, creatorID, userID uuid.UUID) (*domain.Chat, error)
	CreateGroupChat(ctx context.Context, req *domain.CreateGroupChatRequest, creatorID uuid.UUID) (*domain.Chat, error)
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*domain.Chat, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.ChatResponse, error)
	GetAllChats(ctx context.Context) ([]*domain.Chat, error)
	UpdateChat(ctx context.Context, chatID uuid.UUID, req *domain.UpdateChatRequest) (*domain.Chat, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	AddParticipants(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error
	RemoveParticipant(ctx context.Context, chatID, userID, removerID uuid.UUID) error
	LeaveChat(ctx context.Context, chatID, userID uuid.UUID) error
}

type ParticipantService interface {
	AddParticipant(ctx context.Context, chatID, userID, adderID uuid.UUID) error
	AddParticipants(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error
	RemoveParticipant(ctx context.Context, chatID, userID, removerID uuid.UUID) error
	LeaveChat(ctx context.Context, chatID, userID uuid.UUID) error
	GetChatParticipants(ctx context.Context, chatID uuid.UUID) ([]*domain.Participant, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.Participant, error)
	IsParticipant(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	UpdateRole(ctx context.Context, chatID, userID, updaterID uuid.UUID, role string) error
	UpdateLastRead(ctx context.Context, chatID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, chatID, userID uuid.UUID) (int, error)
}

type MessageService interface {
	SendMessage(ctx context.Context, chatID uuid.UUID, req *domain.SendMessageRequest, senderID uuid.UUID) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*domain.Message, error)
	EditMessage(ctx context.Context, messageID uuid.UUID, newText string, userID uuid.UUID) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (*domain.Message, error)
	GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*domain.Message, error)
}
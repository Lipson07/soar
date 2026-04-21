package service

import (
	"context"

	"backend/internal/domain"

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

type SecurityService interface {
	GetUserSettings(ctx context.Context, userID uuid.UUID) (*domain.SecuritySettings, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, req *domain.UpdateSecuritySettingsRequest) error
	SetupTwoFactor(ctx context.Context, userID uuid.UUID, username string) (*domain.TwoFactorSetup, error)
	VerifyAndEnableTwoFactor(ctx context.Context, userID uuid.UUID, code string, secret string) (bool, error)
	DisableTwoFactor(ctx context.Context, userID uuid.UUID) error
	VerifyTwoFactorCode(ctx context.Context, userID uuid.UUID, code string) (bool, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, currentToken string) ([]domain.UserSession, error)
	TerminateSession(ctx context.Context, userID uuid.UUID, sessionID int64, currentToken string) error
	TerminateAllOtherSessions(ctx context.Context, userID uuid.UUID, currentToken string) error
	GenerateSecurityReport(ctx context.Context, userID uuid.UUID) (*domain.SecurityReport, error)
	CreateSession(ctx context.Context, userID uuid.UUID, deviceInfo, deviceType, ipAddress, userAgent string) (*domain.UserSession, error)
}

type FileService interface {
	SaveFile(ctx context.Context, file *domain.File) error
	GetAllFiles(ctx context.Context, userID uuid.UUID) ([]*domain.File, error)
	GetFilesByChat(ctx context.Context, chatID uuid.UUID) ([]*domain.File, error)
	DeleteFile(ctx context.Context, id uuid.UUID) error
}
type CallService interface {
	StartCall(ctx context.Context, chatID uuid.UUID, callerID uuid.UUID, calleeID uuid.UUID, callType domain.CallType) (*domain.Call, error)
	AcceptCall(ctx context.Context, callID uuid.UUID) (*domain.Call, error)
	RejectCall(ctx context.Context, callID uuid.UUID) (*domain.Call, error)
	EndCall(ctx context.Context, callID uuid.UUID) (*domain.Call, error)
	GetCallByID(ctx context.Context, callID uuid.UUID) (*domain.Call, error)
	GetActiveCall(ctx context.Context, chatID uuid.UUID) (*domain.Call, error)
	GetUserCalls(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Call, error)
}

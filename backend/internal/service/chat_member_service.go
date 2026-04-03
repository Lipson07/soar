package service

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/repository"

	"github.com/jmoiron/sqlx"
)

type chatMemberService struct {
	chatMemberRepo repository.ChatMemberRepository
	chatRepo       repository.ChatRepository
	userRepo       repository.UserRepository
	db             *sqlx.DB
}

func NewChatMemberService(
	chatMemberRepo repository.ChatMemberRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	db *sqlx.DB,
) ChatMemberService {
	return &chatMemberService{
		chatMemberRepo: chatMemberRepo,
		chatRepo:       chatRepo,
		userRepo:       userRepo,
		db:             db,
	}
}

func (s *chatMemberService) AddMembers(ctx context.Context, chatID int64, userIDs []int64, currentUserID int64) error {
	for _, userID := range userIDs {
		exists, err := s.chatMemberRepo.Exists(ctx, chatID, userID)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		member := &domain.ChatMember{
			ChatID: chatID,
			UserID: userID,
			Role:   "member",
		}

		if err := s.chatMemberRepo.Create(ctx, member); err != nil {
			return fmt.Errorf("failed to add member %d: %w", userID, err)
		}
	}

	return nil
}

func (s *chatMemberService) AddMember(ctx context.Context, chatID int64, req *domain.AddMemberRequest, currentUserID int64) (*domain.ChatMember, error) {
	exists, err := s.chatMemberRepo.Exists(ctx, chatID, req.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user is already a member of this chat")
	}

	member := &domain.ChatMember{
		ChatID: chatID,
		UserID: req.UserID,
		Role:   "member",
	}

	if err := s.chatMemberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return member, nil
}

func (s *chatMemberService) GetMember(ctx context.Context, chatID, userID int64) (*domain.ChatMember, error) {
	return s.chatMemberRepo.GetByChatAndUser(ctx, chatID, userID)
}

func (s *chatMemberService) GetChatMembers(ctx context.Context, chatID int64) ([]domain.ChatMember, error) {
	return s.chatMemberRepo.GetByChatID(ctx, chatID)
}

func (s *chatMemberService) GetUserChats(ctx context.Context, userID int64) ([]domain.Chat, error) {
	var chats []domain.Chat

	query := `
		SELECT c.* FROM chats c
		INNER JOIN chat_members cm ON c.id = cm.chat_id
		WHERE cm.user_id = $1
		ORDER BY c.updated_at DESC
	`

	err := s.db.SelectContext(ctx, &chats, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}

	return chats, nil
}

func (s *chatMemberService) GetMemberCount(ctx context.Context, chatID int64) (int, error) {
	return s.chatMemberRepo.CountMembers(ctx, chatID)
}

func (s *chatMemberService) UpdateMemberRole(ctx context.Context, chatID, userID int64, req *domain.UpdateMemberRoleRequest, currentUserID int64) error {
	return s.chatMemberRepo.UpdateRole(ctx, chatID, userID, req.Role)
}

func (s *chatMemberService) UpdateLastRead(ctx context.Context, chatID, userID int64) error {
	return s.chatMemberRepo.UpdateLastReadAt(ctx, chatID, userID)
}

func (s *chatMemberService) RemoveMember(ctx context.Context, chatID, userID int64, currentUserID int64) error {
	return s.chatMemberRepo.DeleteByChatAndUser(ctx, chatID, userID)
}

func (s *chatMemberService) LeaveChat(ctx context.Context, chatID, userID int64) error {
	return s.chatMemberRepo.DeleteByChatAndUser(ctx, chatID, userID)
}

func (s *chatMemberService) KickMember(ctx context.Context, chatID, userID int64, currentUserID int64) error {
	return s.RemoveMember(ctx, chatID, userID, currentUserID)
}

func (s *chatMemberService) IsMember(ctx context.Context, chatID, userID int64) (bool, error) {
	return s.chatMemberRepo.Exists(ctx, chatID, userID)
}

func (s *chatMemberService) IsAdmin(ctx context.Context, chatID, userID int64) (bool, error) {
	role, err := s.chatMemberRepo.GetUserRole(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	return role == "admin" || role == "owner", nil
}

func (s *chatMemberService) IsOwner(ctx context.Context, chatID, userID int64) (bool, error) {
	role, err := s.chatMemberRepo.GetUserRole(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	return role == "owner", nil
}

func (s *chatMemberService) GetUserRole(ctx context.Context, chatID, userID int64) (string, error) {
	return s.chatMemberRepo.GetUserRole(ctx, chatID, userID)
}

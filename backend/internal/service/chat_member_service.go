package service

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/repository"
)

type chatMemberService struct {
	chatMemberRepo repository.ChatMemberRepository
	chatRepo       repository.ChatRepository
	userRepo       repository.UserRepository
}

func NewChatMemberService(
	chatMemberRepo repository.ChatMemberRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
) ChatMemberService {
	return &chatMemberService{
		chatMemberRepo: chatMemberRepo,
		chatRepo:       chatRepo,
		userRepo:       userRepo,
	}
}

func (s *chatMemberService) AddMember(ctx context.Context, chatID int64, req *domain.AddMemberRequest, currentUserID int64) (*domain.ChatMember, error) {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}
	if chat == nil {
		return nil, errors.New("chat not found")
	}

	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if chat.Type != "private" {
		isAdmin, err := s.IsAdmin(ctx, chatID, currentUserID)
		if err != nil {
			return nil, err
		}
		if !isAdmin && currentUserID != *chat.CreatedBy {
			return nil, errors.New("permission denied: only admin or owner can add members")
		}
	}

	exists, err := s.chatMemberRepo.Exists(ctx, chatID, req.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user is already a member of this chat")
	}

	role := "member"
	if req.Role != "" {
		role = req.Role
	}

	member := &domain.ChatMember{
		ChatID: chatID,
		UserID: req.UserID,
		Role:   role,
	}

	if err := s.chatMemberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return member, nil
}

func (s *chatMemberService) AddMembers(ctx context.Context, chatID int64, userIDs []int64, currentUserID int64) error {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("chat not found")
	}

	isAdmin, err := s.IsAdmin(ctx, chatID, currentUserID)
	if err != nil {
		return err
	}
	if !isAdmin && currentUserID != *chat.CreatedBy {
		return errors.New("permission denied: only admin or owner can add members")
	}

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

func (s *chatMemberService) GetMember(ctx context.Context, chatID, userID int64) (*domain.ChatMember, error) {
	member, err := s.chatMemberRepo.GetByChatAndUser(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("member not found")
	}
	return member, nil
}

func (s *chatMemberService) GetChatMembers(ctx context.Context, chatID int64) ([]domain.ChatMember, error) {
	return s.chatMemberRepo.GetByChatID(ctx, chatID)
}

func (s *chatMemberService) GetUserChats(ctx context.Context, userID int64) ([]domain.ChatMember, error) {
	return s.chatMemberRepo.GetByUserID(ctx, userID)
}

func (s *chatMemberService) GetMemberCount(ctx context.Context, chatID int64) (int, error) {
	return s.chatMemberRepo.CountMembers(ctx, chatID)
}

func (s *chatMemberService) UpdateMemberRole(ctx context.Context, chatID, userID int64, req *domain.UpdateMemberRoleRequest, currentUserID int64) error {
	isOwner, err := s.IsOwner(ctx, chatID, currentUserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("permission denied: only owner can update roles")
	}

	targetRole, err := s.GetUserRole(ctx, chatID, userID)
	if err != nil {
		return err
	}
	if targetRole == "owner" {
		return errors.New("cannot change owner's role")
	}

	return s.chatMemberRepo.UpdateRole(ctx, chatID, userID, req.Role)
}

func (s *chatMemberService) UpdateLastRead(ctx context.Context, chatID, userID int64) error {
	return s.chatMemberRepo.UpdateLastReadAt(ctx, chatID, userID)
}

func (s *chatMemberService) RemoveMember(ctx context.Context, chatID, userID int64, currentUserID int64) error {
	isAdmin, err := s.IsAdmin(ctx, chatID, currentUserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("permission denied: only admin or owner can remove members")
	}

	if currentUserID == userID {
		return errors.New("cannot remove yourself, use LeaveChat instead")
	}

	targetRole, err := s.GetUserRole(ctx, chatID, userID)
	if err != nil {
		return err
	}
	if targetRole == "owner" {
		return errors.New("cannot remove owner")
	}

	return s.chatMemberRepo.DeleteByChatAndUser(ctx, chatID, userID)
}

func (s *chatMemberService) LeaveChat(ctx context.Context, chatID, userID int64) error {
	isOwner, err := s.IsOwner(ctx, chatID, userID)
	if err != nil {
		return err
	}
	if isOwner {
		return errors.New("owner cannot leave chat, delete it instead")
	}

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

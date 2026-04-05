package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type participantService struct {
	participantRepo repository.ParticipantRepository
	chatRepo        repository.ChatRepository
	userRepo        repository.UserRepository
	messageRepo     repository.MessageRepository
}

func NewParticipantService(
	participantRepo repository.ParticipantRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
) ParticipantService {
	return &participantService{
		participantRepo: participantRepo,
		chatRepo:        chatRepo,
		userRepo:        userRepo,
		messageRepo:     messageRepo,
	}
}

func (s *participantService) AddParticipant(ctx context.Context, chatID, userID, adderID uuid.UUID) error {
	return s.AddParticipants(ctx, chatID, []uuid.UUID{userID}, adderID)
}

func (s *participantService) AddParticipants(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error {
	isAdmin, err := s.participantRepo.IsAdmin(ctx, chatID, adderID)
	if err != nil || !isAdmin {
		return fmt.Errorf("ошибка")
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Type != domain.ChatTypeGroup {
		return domain.ErrInvalidChatType
	}

	for _, uid := range userIDs {
		exists, _ := s.userRepo.Exists(ctx, uid)
		if !exists {
			continue
		}

		isParticipant, _ := s.participantRepo.IsParticipant(ctx, chatID, uid)
		if isParticipant {
			continue
		}

		participant := &domain.Participant{
			ID:         uuid.New(),
			ChatID:     chatID,
			UserID:     uid,
			Role:       string(domain.RoleMember),
			JoinedAt:   time.Now(),
			LastReadAt: time.Now(),
		}
		err = s.participantRepo.Add(ctx, participant)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *participantService) RemoveParticipant(ctx context.Context, chatID, userID, removerID uuid.UUID) error {
	isAdmin, err := s.participantRepo.IsAdmin(ctx, chatID, removerID)
	if err != nil || !isAdmin {
		return fmt.Errorf("ошибка")
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.CreatorID == userID {
		return fmt.Errorf("ошибка")
	}

	return s.participantRepo.Remove(ctx, chatID, userID)
}

func (s *participantService) LeaveChat(ctx context.Context, chatID, userID uuid.UUID) error {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.CreatorID == userID {
		return fmt.Errorf("ошибка")
	}

	return s.participantRepo.Remove(ctx, chatID, userID)
}

func (s *participantService) GetChatParticipants(ctx context.Context, chatID uuid.UUID) ([]*domain.Participant, error) {
	return s.participantRepo.GetByChatID(ctx, chatID)
}

func (s *participantService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.Participant, error) {
	return s.participantRepo.GetByUserID(ctx, userID)
}

func (s *participantService) IsParticipant(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	return s.participantRepo.IsParticipant(ctx, chatID, userID)
}

func (s *participantService) UpdateRole(ctx context.Context, chatID, userID, updaterID uuid.UUID, role string) error {
	isAdmin, err := s.participantRepo.IsAdmin(ctx, chatID, updaterID)
	if err != nil || !isAdmin {
		return fmt.Errorf("ошибка")
	}

	if role != string(domain.RoleAdmin) && role != string(domain.RoleMember) {
		return fmt.Errorf("ошибка")
	}

	participant, err := s.participantRepo.GetByChatAndUser(ctx, chatID, userID)
	if err != nil {
		return err
	}

	participant.Role = role
	return s.participantRepo.Update(ctx, participant)
}

func (s *participantService) UpdateLastRead(ctx context.Context, chatID, userID uuid.UUID) error {
	return s.participantRepo.UpdateLastRead(ctx, chatID, userID, time.Now())
}

func (s *participantService) GetUnreadCount(ctx context.Context, chatID, userID uuid.UUID) (int, error) {
	participant, err := s.participantRepo.GetByChatAndUser(ctx, chatID, userID)
	if err != nil {
		return 0, err
	}
	return s.participantRepo.GetUnreadCount(ctx, chatID, userID, participant.LastReadAt)
}
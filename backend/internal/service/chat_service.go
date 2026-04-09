package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type chatService struct {
	chatRepo        repository.ChatRepository
	participantRepo repository.ParticipantRepository
	userRepo        repository.UserRepository
	messageRepo     repository.MessageRepository
}

func NewChatService(
	chatRepo repository.ChatRepository,
	participantRepo repository.ParticipantRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
) ChatService {
	return &chatService{
		chatRepo:        chatRepo,
		participantRepo: participantRepo,
		userRepo:        userRepo,
		messageRepo:     messageRepo,
	}
}

func (s *chatService) CreatePrivateChat(ctx context.Context, creatorID, userID uuid.UUID) (*domain.Chat, error) {
	if creatorID == userID {
		return nil, errors.New("нельзя создать чат с самим собой")
	}

	// Проверяем существование пользователя
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return nil, errors.New("ошибка проверки пользователя")
	}
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// Проверяем, существует ли уже чат
	existingChat, err := s.chatRepo.GetPrivateChatByUsers(ctx, creatorID, userID)
	if err == nil && existingChat != nil {
		return existingChat, nil
	}

	// Создаем чат БЕЗ имени (nil) - фронтенд сам определит имя
	chat := &domain.Chat{
		ID:        uuid.New(),
		Name:      nil,
		Type:      domain.ChatTypePrivate,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, errors.New("ошибка создания чата")
	}

	// Добавляем участников
	participants := []*domain.Participant{
		{
			ID:         uuid.New(),
			ChatID:     chat.ID,
			UserID:     creatorID,
			Role:       string(domain.RoleMember),
			JoinedAt:   time.Now(),
			LastReadAt: time.Now(),
		},
		{
			ID:         uuid.New(),
			ChatID:     chat.ID,
			UserID:     userID,
			Role:       string(domain.RoleMember),
			JoinedAt:   time.Now(),
			LastReadAt: time.Now(),
		},
	}

	for _, p := range participants {
		if err := s.participantRepo.Add(ctx, p); err != nil {
			return nil, errors.New("ошибка добавления участника")
		}
	}

	return chat, nil
}

func (s *chatService) CreateGroupChat(ctx context.Context, req *domain.CreateGroupChatRequest, creatorID uuid.UUID) (*domain.Chat, error) {
	if req.Name == "" {
		return nil, errors.New("название чата обязательно")
	}

	chat := &domain.Chat{
		ID:        uuid.New(),
		Name:      &req.Name,
		Type:      domain.ChatTypeGroup,
		CreatorID: creatorID,
		AvatarURL: req.AvatarURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, errors.New("ошибка создания группы")
	}

	// Добавляем создателя как админа
	creator := &domain.Participant{
		ID:         uuid.New(),
		ChatID:     chat.ID,
		UserID:     creatorID,
		Role:       string(domain.RoleAdmin),
		JoinedAt:   time.Now(),
		LastReadAt: time.Now(),
	}
	if err := s.participantRepo.Add(ctx, creator); err != nil {
		return nil, errors.New("ошибка добавления создателя")
	}

	// Добавляем остальных участников
	for _, userID := range req.UserIDs {
		if userID == creatorID {
			continue
		}
		participant := &domain.Participant{
			ID:         uuid.New(),
			ChatID:     chat.ID,
			UserID:     userID,
			Role:       string(domain.RoleMember),
			JoinedAt:   time.Now(),
			LastReadAt: time.Now(),
		}
		if err := s.participantRepo.Add(ctx, participant); err != nil {
			return nil, errors.New("ошибка добавления участника")
		}
	}

	return chat, nil
}

func (s *chatService) GetChatByID(ctx context.Context, chatID uuid.UUID) (*domain.Chat, error) {
	return s.chatRepo.GetByID(ctx, chatID)
}

func (s *chatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.ChatResponse, error) {
	participants, err := s.participantRepo.GetByUserID(ctx, userID)
	if err != nil {
		return []*domain.ChatResponse{}, nil
	}

	var chats []*domain.ChatResponse
	for _, p := range participants {
		chat, err := s.chatRepo.GetByID(ctx, p.ChatID)
		if err != nil {
			continue
		}

		lastMessage, _ := s.messageRepo.GetLastMessage(ctx, chat.ID)
		unreadCount, _ := s.participantRepo.GetUnreadCount(ctx, chat.ID, userID, p.LastReadAt)

		// Для приватных чатов, если имя не установлено, получаем имя другого участника
		if chat.Type == domain.ChatTypePrivate && (chat.Name == nil || *chat.Name == "") {
			participantsInChat, _ := s.participantRepo.GetByChatID(ctx, chat.ID)
			for _, part := range participantsInChat {
				if part.UserID != userID {
					otherUser, _ := s.userRepo.GetByID(ctx, part.UserID)
					if otherUser != nil {
						chat.Name = &otherUser.Username
					}
					break
				}
			}
		}

		chatResp := &domain.ChatResponse{
			ID:            chat.ID,
			Name:          chat.Name,
			Type:          chat.Type,
			CreatorID:     chat.CreatorID,
			AvatarURL:     chat.AvatarURL,
			LastMessage:   lastMessage,
			UnreadCount:   unreadCount,
			LastMessageAt: chat.LastMessageAt,
			UpdatedAt:     chat.UpdatedAt,
		}
		chats = append(chats, chatResp)
	}

	return chats, nil
}

func (s *chatService) GetAllChats(ctx context.Context) ([]*domain.Chat, error) {
	return s.chatRepo.GetAll(ctx)
}

func (s *chatService) UpdateChat(ctx context.Context, chatID uuid.UUID, req *domain.UpdateChatRequest) (*domain.Chat, error) {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, errors.New("чат не найден")
	}

	if req.Name != nil {
		chat.Name = req.Name
	}
	if req.AvatarURL != nil {
		chat.AvatarURL = req.AvatarURL
	}
	chat.UpdatedAt = time.Now()

	if err := s.chatRepo.Update(ctx, chat); err != nil {
		return nil, errors.New("ошибка обновления чата")
	}

	return chat, nil
}

func (s *chatService) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	return s.chatRepo.Delete(ctx, chatID)
}

func (s *chatService) AddParticipants(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error {
	isAdmin, err := s.participantRepo.IsAdmin(ctx, chatID, adderID)
	if err != nil || !isAdmin {
		return errors.New("только администраторы могут добавлять участников")
	}

	for _, userID := range userIDs {
		participant := &domain.Participant{
			ID:         uuid.New(),
			ChatID:     chatID,
			UserID:     userID,
			Role:       string(domain.RoleMember),
			JoinedAt:   time.Now(),
			LastReadAt: time.Now(),
		}
		if err := s.participantRepo.Add(ctx, participant); err != nil {
			return err
		}
	}
	return nil
}

func (s *chatService) RemoveParticipant(ctx context.Context, chatID, userID, removerID uuid.UUID) error {
	isAdmin, err := s.participantRepo.IsAdmin(ctx, chatID, removerID)
	if err != nil || !isAdmin {
		return errors.New("только администраторы могут удалять участников")
	}
	return s.participantRepo.Remove(ctx, chatID, userID)
}

func (s *chatService) LeaveChat(ctx context.Context, chatID, userID uuid.UUID) error {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat.CreatorID == userID {
		return errors.New("создатель не может покинуть чат")
	}
	return s.participantRepo.Remove(ctx, chatID, userID)
}

package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type messageService struct {
	messageRepo     repository.MessageRepository
	participantRepo repository.ParticipantRepository
	chatRepo        repository.ChatRepository
}

func NewMessageService(
	messageRepo repository.MessageRepository,
	participantRepo repository.ParticipantRepository,
	chatRepo repository.ChatRepository,
) MessageService {
	return &messageService{
		messageRepo:     messageRepo,
		participantRepo: participantRepo,
		chatRepo:        chatRepo,
	}
}

func (s *messageService) SendMessage(ctx context.Context, chatID uuid.UUID, req *domain.SendMessageRequest, senderID uuid.UUID) (*domain.Message, error) {
	isParticipant, err := s.participantRepo.IsParticipant(ctx, chatID, senderID)
	if err != nil || !isParticipant {
		return nil, fmt.Errorf("пользователь не является участником чата")
	}

	// Валидация в зависимости от типа сообщения
	switch req.Type {
	case domain.MessageTypeText:
		if req.Text == "" {
			return nil, fmt.Errorf("текст сообщения не может быть пустым")
		}
	case domain.MessageTypeImage, domain.MessageTypeFile:
		if req.FileURL == "" {
			return nil, fmt.Errorf("URL файла обязателен")
		}
	default:
		req.Type = domain.MessageTypeText
	}

	message := &domain.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		UserID:    senderID,
		Type:      req.Type,
		Text:      req.Text,
		FileURL:   req.FileURL,
		FileName:  req.FileName,
		FileSize:  req.FileSize,
		MimeType:  req.MimeType,
		ReplyTo:   req.ReplyTo,
		IsEdited:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	s.chatRepo.UpdateLastMessage(ctx, chatID, message.CreatedAt)

	return message, nil
}

func (s *messageService) GetMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	isParticipant, err := s.participantRepo.IsParticipant(ctx, chatID, userID)
	if err != nil || !isParticipant {
		return nil, fmt.Errorf("пользователь не является участником чата")
	}

	return s.messageRepo.GetByChatID(ctx, chatID, limit, offset)
}

func (s *messageService) EditMessage(ctx context.Context, messageID uuid.UUID, newText string, userID uuid.UUID) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return fmt.Errorf("сообщение не найдено")
	}

	// Редактировать можно только текстовые сообщения
	if message.Type != domain.MessageTypeText {
		return fmt.Errorf("можно редактировать только текстовые сообщения")
	}

	if message.UserID != userID {
		return fmt.Errorf("нет прав на редактирование сообщения")
	}

	message.Text = newText
	message.IsEdited = true
	message.UpdatedAt = time.Now()

	return s.messageRepo.Update(ctx, message)
}

func (s *messageService) DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return fmt.Errorf("сообщение не найдено")
	}

	if message.UserID == userID {
		return s.messageRepo.Delete(ctx, messageID)
	}

	isAdmin, err := s.participantRepo.IsAdmin(ctx, message.ChatID, userID)
	if err != nil || !isAdmin {
		return fmt.Errorf("нет прав на удаление сообщения")
	}

	return s.messageRepo.Delete(ctx, messageID)
}

func (s *messageService) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*domain.Message, error) {
	return s.messageRepo.GetByID(ctx, messageID)
}

func (s *messageService) GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	return s.messageRepo.GetByChatID(ctx, chatID, limit, offset)
}
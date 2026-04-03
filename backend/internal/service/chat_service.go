package service

import (
	"context"
	"fmt"

	"myapp/internal/domain"
	"myapp/internal/repository"
)

type chatService struct {
	chatRepo       repository.ChatRepository
	chatMemberRepo repository.ChatMemberRepository
}

func NewChatService(
	chatRepo repository.ChatRepository,
	chatMemberRepo repository.ChatMemberRepository,
) ChatService {
	return &chatService{
		chatRepo:       chatRepo,
		chatMemberRepo: chatMemberRepo,
	}
}

func (s *chatService) Create(ctx context.Context, req *domain.CreateChatRequest, userID int64) (*domain.Chat, error) {
	if req.Type != "private" && (req.Name == nil || *req.Name == "") {
		return nil, fmt.Errorf("название чата не может быть пустым для групповых чатов")
	}

	if req.Name != nil && *req.Name != "" {
		existing, _ := s.chatRepo.GetByName(ctx, *req.Name)
		if existing != nil {
			return nil, domain.ErrChatNameExists
		}
	}

	var name string
	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	} else if req.Type == "private" {
		name = ""
	}

	chat := &domain.Chat{
		Type:        req.Type,
		Name:        &name,
		Description: req.Description,
		AvatarPath:  nil,
		CreatedBy:   &userID,
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	member := &domain.ChatMember{
		ChatID: chat.ID,
		UserID: userID,
		Role:   "owner",
	}

	if err := s.chatMemberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) GetByID(ctx context.Context, id int64) (*domain.Chat, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidChatID
	}

	chat, err := s.chatRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) GetByName(ctx context.Context, name string) (*domain.Chat, error) {
	if name == "" {
		return nil, fmt.Errorf("название чата не может быть пустым")
	}

	chat, err := s.chatRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) GetAll(ctx context.Context) ([]domain.Chat, error) {
	chats, err := s.chatRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (s *chatService) Update(ctx context.Context, id int64, req *domain.UpdateChatRequest) (*domain.Chat, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidChatID
	}

	chat, err := s.chatRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if chat.Type != "private" && *req.Name == "" {
			return nil, fmt.Errorf("название чата не может быть пустым")
		}

		currentName := ""
		if chat.Name != nil {
			currentName = *chat.Name
		}

		if *req.Name != currentName && *req.Name != "" {
			existing, _ := s.chatRepo.GetByName(ctx, *req.Name)
			if existing != nil {
				return nil, domain.ErrChatNameExists
			}
		}

		if chat.Type != "private" || *req.Name != "" {
			chat.Name = req.Name
		}
	}

	if req.Description != nil {
		chat.Description = req.Description
	}

	if req.AvatarPath != nil {
		chat.AvatarPath = req.AvatarPath
	}

	if err := s.chatRepo.Update(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidChatID
	}
	return s.chatRepo.Delete(ctx, id)
}

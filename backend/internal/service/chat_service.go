package service

import (
	"context"
	"fmt"

	"myapp/internal/domain"
	"myapp/internal/repository"
)

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) Create(ctx context.Context, req *domain.CreateChatRequest) (*domain.Chat, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("название чата не может быть пустым")
	}

	existing, _ := s.chatRepo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, domain.ErrChatNameExists
	}

	chat := &domain.Chat{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		AvatarPath:  nil,
		CreatedBy:   nil,
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
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
		if *req.Name == "" {
			return nil, fmt.Errorf("название чата не может быть пустым")
		}
		if *req.Name != chat.Name {
			existing, _ := s.chatRepo.GetByName(ctx, *req.Name)
			if existing != nil {
				return nil, domain.ErrChatNameExists
			}
		}
		chat.Name = *req.Name
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
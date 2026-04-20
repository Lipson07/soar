package service

import (
	"context"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type fileService struct {
	fileRepo repository.FileRepository
}

func NewFileService(fileRepo repository.FileRepository) FileService {
	return &fileService{
		fileRepo: fileRepo,
	}
}

func (s *fileService) SaveFile(ctx context.Context, file *domain.File) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	return s.fileRepo.Create(ctx, file)
}

func (s *fileService) GetAllFiles(ctx context.Context, userID uuid.UUID) ([]*domain.File, error) {
	return s.fileRepo.GetByUserID(ctx, userID)
}

func (s *fileService) GetFilesByChat(ctx context.Context, chatID uuid.UUID) ([]*domain.File, error) {
	return s.fileRepo.GetByChatID(ctx, chatID)
}

func (s *fileService) DeleteFile(ctx context.Context, id uuid.UUID) error {
	return s.fileRepo.Delete(ctx, id)
}

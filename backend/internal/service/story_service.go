package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type Storyservice struct {
	repo         repository.StoryRepository
	uploadDir    string
	maxSize      int64
	allowedTypes []string
}

func NewStoryService(repo repository.StoryRepository, uploadDir string) StoryService {
	return &Storyservice{
		repo:      repo,
		uploadDir: uploadDir,
		maxSize:   50 * 1024 * 1024,
		allowedTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"video/mp4",
			"video/webm",
		},
	}
}

func (s *Storyservice) UploadStory(userID, userName string, userAvatar *string, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	if header.Size > s.maxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxSize)
	}

	contentType := header.Header.Get("Content-Type")
	if !s.isAllowedType(contentType) {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	fileType := "image"
	if strings.HasPrefix(contentType, "video/") {
		fileType = "video"
	}

	userDir := filepath.Join(s.uploadDir, "stories", userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(userDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	fileURL := fmt.Sprintf("/uploads/stories/%s/%s", userID, fileName)

	story := &domain.Story{
		UserID:     userID,
		UserName:   userName,
		UserAvatar: userAvatar,
		FileURL:    fileURL,
		Type:       fileType,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	createdStory, err := s.repo.Create(story)
	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create story record: %w", err)
	}

	return createdStory, nil
}

func (s *Storyservice) GetStories(userID string) ([]domain.Story, error) {
	stories, err := s.repo.GetAll(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stories: %w", err)
	}

	return stories, nil
}

func (s *Storyservice) MarkStoryAsViewed(storyID, userID string) error {
	return s.repo.MarkAsViewed(storyID, userID)
}

func (s *Storyservice) isAllowedType(contentType string) bool {
	for _, allowedType := range s.allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/jmoiron/sqlx"
)

type StoryRepo struct {
	db *sqlx.DB
}

func NewStoryRepo(db *sqlx.DB) repository.StoryRepository {
	return &StoryRepo{db: db}
}

func (r *StoryRepo) Create(story *domain.Story) (*domain.Story, error) {
	query := `
		INSERT INTO stories (user_id, user_name, user_avatar, file_url, type, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, expires_at`

	now := time.Now()
	story.CreatedAt = now
	story.ExpiresAt = now.Add(24 * time.Hour)

	err := r.db.QueryRow(
		query,
		story.UserID,
		story.UserName,
		story.UserAvatar,
		story.FileURL,
		story.Type,
		story.CreatedAt,
		story.ExpiresAt,
	).Scan(&story.ID, &story.CreatedAt, &story.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create story: %w", err)
	}

	return story, nil
}

func (r *StoryRepo) GetAll(userID string) ([]domain.Story, error) {
	query := `
		SELECT 
			s.id, 
			s.user_id, 
			s.user_name, 
			s.user_avatar, 
			s.file_url, 
			s.type, 
			s.created_at, 
			s.expires_at,
			CASE WHEN sv.story_id IS NOT NULL THEN true ELSE false END as viewed
		FROM stories s
		LEFT JOIN story_views sv ON s.id = sv.story_id AND sv.user_id = $1
		WHERE s.expires_at > NOW()
		ORDER BY 
			CASE WHEN s.user_id = $1 THEN 0 ELSE 1 END,
			s.created_at DESC`

	var stories []domain.Story
	err := r.db.Select(&stories, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stories: %w", err)
	}

	return stories, nil
}

func (r *StoryRepo) GetByID(id string) (*domain.Story, error) {
	query := `
		SELECT id, user_id, user_name, user_avatar, file_url, type, created_at, expires_at
		FROM stories
		WHERE id = $1 AND expires_at > NOW()`

	var story domain.Story
	err := r.db.Get(&story, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("story not found")
		}
		return nil, fmt.Errorf("failed to get story: %w", err)
	}

	return &story, nil
}

func (r *StoryRepo) DeleteExpired() error {
	query := `DELETE FROM stories WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired stories: %w", err)
	}
	return nil
}

func (r *StoryRepo) MarkAsViewed(storyID, userID string) error {
	query := `
		INSERT INTO story_views (story_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (story_id, user_id) DO NOTHING`

	_, err := r.db.Exec(query, storyID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark story as viewed: %w", err)
	}
	return nil
}

func (r *StoryRepo) GetViewers(storyID string) ([]string, error) {
	query := `SELECT user_id FROM story_views WHERE story_id = $1 ORDER BY viewed_at DESC`

	var viewers []string
	err := r.db.Select(&viewers, query, storyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get viewers: %w", err)
	}

	return viewers, nil
}

func (r *StoryRepo) StartExpiredStoriesCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := r.DeleteExpired(); err != nil {
				fmt.Printf("Error cleaning expired stories: %v\n", err)
			}
		}
	}()
}

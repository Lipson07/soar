package postgres

import (
	"context"
	"database/sql"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) repository.FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, file *domain.File) error {
	query := `
		INSERT INTO files (id, name, url, size, type, mime_type, chat_id, message_id, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		file.ID, file.Name, file.URL, file.Size, file.Type,
		file.MimeType, file.ChatID, file.MessageID, file.UserID, file.CreatedAt,
	)
	return err
}

func (r *FileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	query := `SELECT id, name, url, size, type, mime_type, chat_id, message_id, user_id, created_at FROM files WHERE id = $1`
	var file domain.File
	err := r.db.GetContext(ctx, &file, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &file, err
}

func (r *FileRepository) GetAll(ctx context.Context) ([]*domain.File, error) {
	query := `SELECT id, name, url, size, type, mime_type, chat_id, message_id, user_id, created_at FROM files ORDER BY created_at DESC`
	var files []*domain.File
	err := r.db.SelectContext(ctx, &files, query)
	return files, err
}

func (r *FileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.File, error) {
	query := `SELECT id, name, url, size, type, mime_type, chat_id, message_id, user_id, created_at FROM files WHERE user_id = $1 ORDER BY created_at DESC`
	var files []*domain.File
	err := r.db.SelectContext(ctx, &files, query, userID)
	return files, err
}

func (r *FileRepository) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]*domain.File, error) {
	query := `SELECT id, name, url, size, type, mime_type, chat_id, message_id, user_id, created_at FROM files WHERE chat_id = $1 ORDER BY created_at DESC`
	var files []*domain.File
	err := r.db.SelectContext(ctx, &files, query, chatID)
	return files, err
}

func (r *FileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

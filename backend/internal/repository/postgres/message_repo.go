package postgres

import (
	"backend/internal/domain"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO messages (id, chat_id, sender_id, text, reply_to, is_edited, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		message.ID, message.ChatID, message.UserID,
		message.Text, message.ReplyTo, message.IsEdited,
		message.CreatedAt, message.UpdatedAt,
	)
	return err
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	query := `SELECT id, chat_id, sender_id, text, reply_to, is_edited, created_at, updated_at, deleted_at 
		FROM messages WHERE id = $1 AND deleted_at IS NULL`
	var message domain.Message
	err := r.db.GetContext(ctx, &message, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &message, err
}

func (r *MessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	query := `
		SELECT id, chat_id, sender_id, text, reply_to, is_edited, created_at, updated_at, deleted_at
		FROM messages
		WHERE chat_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var messages []*domain.Message
	err := r.db.SelectContext(ctx, &messages, query, chatID, limit, offset)
	return messages, err
}

func (r *MessageRepository) Update(ctx context.Context, message *domain.Message) error {
	query := `
		UPDATE messages SET text = $1, is_edited = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, message.Text, message.IsEdited, message.UpdatedAt, message.ID)
	return err
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE messages SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *MessageRepository) GetLastMessage(ctx context.Context, chatID uuid.UUID) (*domain.MessageInfo, error) {
	query := `
		SELECT id, text, sender_id, created_at
		FROM messages
		WHERE chat_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	var msg domain.MessageInfo
	err := r.db.GetContext(ctx, &msg, query, chatID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &msg, err
}
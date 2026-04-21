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
		INSERT INTO messages (id, chat_id, user_id, type, text, file_url, file_name, file_size, mime_type, reply_to, is_edited, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		message.ID, message.ChatID, message.UserID, message.Type,
		message.Text, message.FileURL, message.FileName, message.FileSize, message.MimeType,
		message.ReplyTo, message.IsEdited, message.CreatedAt, message.UpdatedAt,
	)
	return err
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	query := `SELECT id, chat_id, user_id, type, text, file_url, file_name, file_size, mime_type, reply_to, is_edited, created_at, updated_at, deleted_at 
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
		SELECT id, chat_id, user_id, type, text, file_url, file_name, file_size, mime_type, reply_to, is_edited, created_at, updated_at, deleted_at
		FROM messages
		WHERE chat_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
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
	_, err := r.db.ExecContext(ctx, query, message.Text, true, message.UpdatedAt, message.ID)
	return err
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE messages SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
func (r *MessageRepository) GetLastMessage(ctx context.Context, chatID uuid.UUID) (*domain.MessageInfo, error) {
	query := `
		SELECT id, type, text, file_url, file_name, user_id, created_at
		FROM messages
		WHERE chat_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	var msg domain.MessageInfo
	err := r.db.GetContext(ctx, &msg, query, chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if msg.ID == uuid.Nil {
		return nil, nil
	}

	return &msg, nil
}
func (r *MessageRepository) UpdateStatus(ctx context.Context, messageID uuid.UUID, status domain.MessageStatus) error {
	return nil
}

func (r *MessageRepository) AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	query := `
		INSERT INTO reactions (message_id, user_id, emoji, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (message_id, user_id, emoji) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, messageID, userID, emoji)
	return err
}

func (r *MessageRepository) RemoveReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	query := `DELETE FROM reactions WHERE message_id = $1 AND user_id = $2 AND emoji = $3`
	_, err := r.db.ExecContext(ctx, query, messageID, userID, emoji)
	return err
}

func (r *MessageRepository) GetReactions(ctx context.Context, messageID uuid.UUID) ([]domain.Reaction, error) {
	query := `SELECT message_id, user_id, emoji, created_at FROM reactions WHERE message_id = $1`
	var reactions []domain.Reaction
	err := r.db.SelectContext(ctx, &reactions, query, messageID)
	return reactions, err
}

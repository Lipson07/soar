package postgres

import (
	"backend/internal/domain"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	query := `
		INSERT INTO chats (id, name, type, creator_id, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		chat.ID, chat.Name, chat.Type, chat.CreatorID, chat.AvatarURL, chat.CreatedAt, chat.UpdatedAt,
	)
	return err
}

func (r *ChatRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	query := `SELECT id, name, type, creator_id, avatar_url, created_at, updated_at, last_message_at FROM chats WHERE id = $1`
	var chat domain.Chat
	err := r.db.GetContext(ctx, &chat, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrChatNotFound
	}
	return &chat, err
}

func (r *ChatRepository) GetPrivateChatByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*domain.Chat, error) {
	query := `
		SELECT c.id, c.name, c.type, c.creator_id, c.avatar_url, c.created_at, c.updated_at, c.last_message_at
		FROM chats c
		JOIN participants p1 ON p1.chat_id = c.id
		JOIN participants p2 ON p2.chat_id = c.id
		WHERE c.type = $1 
		AND p1.user_id = $2 
		AND p2.user_id = $3
		LIMIT 1
	`
	var chat domain.Chat
	err := r.db.GetContext(ctx, &chat, query, domain.ChatTypePrivate, user1ID, user2ID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &chat, err
}

func (r *ChatRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	query := `
		SELECT c.id, c.name, c.type, c.creator_id, c.avatar_url, c.created_at, c.updated_at, c.last_message_at
		FROM chats c
		JOIN participants p ON p.chat_id = c.id
		WHERE p.user_id = $1
		ORDER BY c.updated_at DESC
	`
	var chats []*domain.Chat
	err := r.db.SelectContext(ctx, &chats, query, userID)
	return chats, err
}

func (r *ChatRepository) GetAll(ctx context.Context) ([]*domain.Chat, error) {
	query := `SELECT id, name, type, creator_id, avatar_url, created_at, updated_at, last_message_at FROM chats ORDER BY created_at DESC`
	var chats []*domain.Chat
	err := r.db.SelectContext(ctx, &chats, query)
	return chats, err
}

func (r *ChatRepository) Update(ctx context.Context, chat *domain.Chat) error {
	query := `
		UPDATE chats SET name = $1, avatar_url = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, chat.Name, chat.AvatarURL, chat.UpdatedAt, chat.ID)
	return err
}

func (r *ChatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM chats WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ChatRepository) UpdateLastMessage(ctx context.Context, chatID uuid.UUID, lastMessageAt time.Time) error {
	query := `UPDATE chats SET last_message_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, lastMessageAt, chatID)
	return err
}

package postgres

import (
	"context"
	"fmt"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ChatRepository struct {
	db *sqlx.DB
}
func NewChatRepostory(db *sqlx.DB) *ChatRepository{
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context,chat *domain.Chat)error{
	query := `
        INSERT INTO users (type,name,description,avatar_path,created_by,created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
err := r.db.QueryRowContext(
		ctx, query,
		chat.Type,chat.Name,chat.Description,chat.AvatarPath,chat.CreatedBy,
	).Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания: %w", err)
	}
	return nil

}
func (r *ChatRepository) GetByID(ctx context.Context, id int64) (*domain.Chat, error) {
    query := `
        SELECT id, type, name, description, avatar_path, created_by, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `
    
    var chat domain.Chat
    err := r.db.GetContext(ctx, &chat, query, id)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения чата по ID %d: %w", id, err)
    }
    
    return &chat, nil
}

func (r *ChatRepository) GetByName(ctx context.Context, name string) (*domain.Chat, error) {
    query := `
        SELECT id, type, name, description, avatar_path, created_by, created_at, updated_at 
        FROM users 
        WHERE name = $1
    `
    
    var chat domain.Chat
    err := r.db.GetContext(ctx, &chat, query, name)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения чата по имени %s: %w", name, err)
    }
    
    return &chat, nil
}

func (r *ChatRepository) GetAll(ctx context.Context) ([]domain.Chat, error) {
    query := `
        SELECT id, type, name, description, avatar_path, created_by, created_at, updated_at 
        FROM users 
        ORDER BY id
    `
    
    var chats []domain.Chat
    err := r.db.SelectContext(ctx, &chats, query)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения всех чатов: %w", err)
    }
    
    return chats, nil
}

func (r *ChatRepository) Update(ctx context.Context, chat *domain.Chat) error {
    query := `
        UPDATE users 
        SET type = $1, name = $2, description = $3, avatar_path = $4, created_by = $5, updated_at = NOW()
        WHERE id = $6
        RETURNING updated_at
    `
    
    err := r.db.QueryRowContext(
        ctx, query,
        chat.Type, chat.Name, chat.Description, chat.AvatarPath, chat.CreatedBy, chat.ID,
    ).Scan(&chat.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("ошибка обновления чата с ID %d: %w", chat.ID, err)
    }
    
    return nil
}

func (r *ChatRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM users WHERE id = $1`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("ошибка удаления чата с ID %d: %w", id, err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("чат с ID %d не найден", id)
    }
    
    return nil
}
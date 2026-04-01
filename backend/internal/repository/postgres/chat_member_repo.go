package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type ChatMemberRepository struct {
	db *sqlx.DB
}

func NewChatMemberRepository(db *sqlx.DB) *ChatMemberRepository {
	return &ChatMemberRepository{db: db}
}

func (r *ChatMemberRepository) Create(ctx context.Context, member *domain.ChatMember) error {
	query := `
		INSERT INTO chat_members (chat_id, user_id, role, joined_at, last_read_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, joined_at, last_read_at
	`

	now := time.Now()
	member.JoinedAt = now
	member.LastReadAt = now

	err := r.db.QueryRowxContext(ctx, query,
		member.ChatID,
		member.UserID,
		member.Role,
		member.JoinedAt,
		member.LastReadAt,
	).Scan(&member.ID, &member.JoinedAt, &member.LastReadAt)

	if err != nil {
		return fmt.Errorf("failed to create chat member: %w", err)
	}

	return nil
}

func (r *ChatMemberRepository) GetByID(ctx context.Context, id int64) (*domain.ChatMember, error) {
	var member domain.ChatMember
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM chat_members WHERE id = $1`

	err := r.db.GetContext(ctx, &member, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chat member by id: %w", err)
	}

	return &member, nil
}

func (r *ChatMemberRepository) GetByChatAndUser(ctx context.Context, chatID, userID int64) (*domain.ChatMember, error) {
	var member domain.ChatMember
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM chat_members WHERE chat_id = $1 AND user_id = $2`

	err := r.db.GetContext(ctx, &member, query, chatID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chat member: %w", err)
	}

	return &member, nil
}

func (r *ChatMemberRepository) GetByChatID(ctx context.Context, chatID int64) ([]domain.ChatMember, error) {
	var members []domain.ChatMember
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM chat_members WHERE chat_id = $1 ORDER BY joined_at ASC`

	err := r.db.SelectContext(ctx, &members, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat members: %w", err)
	}

	return members, nil
}

func (r *ChatMemberRepository) GetByUserID(ctx context.Context, userID int64) ([]domain.ChatMember, error) {
	var members []domain.ChatMember
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM chat_members WHERE user_id = $1 ORDER BY joined_at DESC`

	err := r.db.SelectContext(ctx, &members, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}

	return members, nil
}

func (r *ChatMemberRepository) Update(ctx context.Context, member *domain.ChatMember) error {
	query := `
		UPDATE chat_members 
		SET role = $1, last_read_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, member.Role, member.LastReadAt, member.ID)
	if err != nil {
		return fmt.Errorf("failed to update chat member: %w", err)
	}

	return nil
}

func (r *ChatMemberRepository) UpdateRole(ctx context.Context, chatID, userID int64, role string) error {
	query := `UPDATE chat_members SET role = $1 WHERE chat_id = $2 AND user_id = $3`

	result, err := r.db.ExecContext(ctx, query, role, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *ChatMemberRepository) UpdateLastReadAt(ctx context.Context, chatID, userID int64) error {
	query := `UPDATE chat_members SET last_read_at = $1 WHERE chat_id = $2 AND user_id = $3`

	_, err := r.db.ExecContext(ctx, query, time.Now(), chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to update last_read_at: %w", err)
	}

	return nil
}

func (r *ChatMemberRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM chat_members WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat member: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *ChatMemberRepository) DeleteByChatAndUser(ctx context.Context, chatID, userID int64) error {
	query := `DELETE FROM chat_members WHERE chat_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete chat member: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *ChatMemberRepository) DeleteAllByChatID(ctx context.Context, chatID int64) error {
	query := `DELETE FROM chat_members WHERE chat_id = $1`

	_, err := r.db.ExecContext(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete all chat members: %w", err)
	}

	return nil
}

func (r *ChatMemberRepository) Exists(ctx context.Context, chatID, userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)`

	err := r.db.QueryRowContext(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check member existence: %w", err)
	}

	return exists, nil
}

func (r *ChatMemberRepository) CountMembers(ctx context.Context, chatID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM chat_members WHERE chat_id = $1`

	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count members: %w", err)
	}

	return count, nil
}

func (r *ChatMemberRepository) GetUserRole(ctx context.Context, chatID, userID int64) (string, error) {
	var role string
	query := `SELECT role FROM chat_members WHERE chat_id = $1 AND user_id = $2`

	err := r.db.QueryRowContext(ctx, query, chatID, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	return role, nil
}

func (r *ChatMemberRepository) IsMember(ctx context.Context, chatID, userID int64) (bool, error) {
	return r.Exists(ctx, chatID, userID)
}

func (r *ChatMemberRepository) IsOwner(ctx context.Context, chatID, userID int64) (bool, error) {
	var isOwner bool
	query := `SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2 AND role = 'owner')`

	err := r.db.QueryRowContext(ctx, query, chatID, userID).Scan(&isOwner)
	if err != nil {
		return false, fmt.Errorf("failed to check owner: %w", err)
	}

	return isOwner, nil
}

func (r *ChatMemberRepository) IsAdmin(ctx context.Context, chatID, userID int64) (bool, error) {
	var isAdmin bool
	query := `SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2 AND role IN ('owner', 'admin'))`

	err := r.db.QueryRowContext(ctx, query, chatID, userID).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("failed to check admin: %w", err)
	}

	return isAdmin, nil
}

package postgres

import (
	"backend/internal/domain"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type participantRepository struct {
	db *sqlx.DB
}

func NewParticipantRepository(db *sqlx.DB) *participantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) Add(ctx context.Context, participant *domain.Participant) error {
	query := `
		INSERT INTO participants (id, chat_id, user_id, role, joined_at, last_read_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query,
		participant.ID, participant.ChatID, participant.UserID,
		participant.Role, participant.JoinedAt, participant.LastReadAt,
	)
	return err
}

func (r *participantRepository) Remove(ctx context.Context, chatID, userID uuid.UUID) error {
	query := `DELETE FROM participants WHERE chat_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, chatID, userID)
	return err
}

func (r *participantRepository) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]*domain.Participant, error) {
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM participants WHERE chat_id = $1`
	var participants []*domain.Participant
	err := r.db.SelectContext(ctx, &participants, query, chatID)
	return participants, err
}

func (r *participantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Participant, error) {
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM participants WHERE user_id = $1`
	var participants []*domain.Participant
	err := r.db.SelectContext(ctx, &participants, query, userID)
	return participants, err
}

func (r *participantRepository) IsParticipant(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM participants WHERE chat_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, chatID, userID)
	return exists, err
}

func (r *participantRepository) IsAdmin(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM participants WHERE chat_id = $1 AND user_id = $2 AND role = $3)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, chatID, userID, domain.RoleAdmin)
	return exists, err
}

func (r *participantRepository) GetByChatAndUser(ctx context.Context, chatID, userID uuid.UUID) (*domain.Participant, error) {
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM participants WHERE chat_id = $1 AND user_id = $2`
	var p domain.Participant
	err := r.db.GetContext(ctx, &p, query, chatID, userID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &p, err
}

func (r *participantRepository) Update(ctx context.Context, participant *domain.Participant) error {
	query := `UPDATE participants SET role = $1, last_read_at = $2 WHERE chat_id = $3 AND user_id = $4`
	_, err := r.db.ExecContext(ctx, query, participant.Role, participant.LastReadAt, participant.ChatID, participant.UserID)
	return err
}

func (r *participantRepository) UpdateLastRead(ctx context.Context, chatID, userID uuid.UUID, lastReadAt time.Time) error {
	query := `UPDATE participants SET last_read_at = $1 WHERE chat_id = $2 AND user_id = $3`
	_, err := r.db.ExecContext(ctx, query, lastReadAt, chatID, userID)
	return err
}

func (r *participantRepository) GetUnreadCount(ctx context.Context, chatID, userID uuid.UUID, lastReadAt time.Time) (int, error) {
	query := `
		SELECT COUNT(*) FROM messages 
		WHERE chat_id = $1 AND created_at > $2 AND deleted_at IS NULL
	`
	var count int
	err := r.db.GetContext(ctx, &count, query, chatID, lastReadAt)
	return count, err
}
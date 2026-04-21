package postgres

import (
	"backend/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ParticipantRepository struct {
	db *sqlx.DB
}

func NewParticipantRepository(db *sqlx.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) Add(ctx context.Context, participant *domain.Participant) error {
	query := `
		INSERT INTO participants (id, chat_id, user_id, role, joined_at, last_read_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		participant.ID, participant.ChatID, participant.UserID,
		participant.Role, participant.JoinedAt, participant.LastReadAt,
	)
	return err
}

func (r *ParticipantRepository) Remove(ctx context.Context, chatID, userID uuid.UUID) error {
	query := `DELETE FROM participants WHERE chat_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, chatID, userID)
	return err
}

func (r *ParticipantRepository) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]*domain.Participant, error) {
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM participants WHERE chat_id = $1`
	var participants []*domain.Participant
	err := r.db.SelectContext(ctx, &participants, query, chatID)
	return participants, err
}

func (r *ParticipantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Participant, error) {
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM participants WHERE user_id = $1`
	var participants []*domain.Participant
	err := r.db.SelectContext(ctx, &participants, query, userID)
	return participants, err
}

func (r *ParticipantRepository) IsParticipant(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM participants WHERE chat_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, chatID, userID)
	return exists, err
}

func (r *ParticipantRepository) IsAdmin(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM participants WHERE chat_id = $1 AND user_id = $2 AND role = 'admin')`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, chatID, userID)
	return exists, err
}

func (r *ParticipantRepository) GetByChatAndUser(ctx context.Context, chatID, userID uuid.UUID) (*domain.Participant, error) {
	query := `SELECT id, chat_id, user_id, role, joined_at, last_read_at FROM participants WHERE chat_id = $1 AND user_id = $2`
	var participant domain.Participant
	err := r.db.GetContext(ctx, &participant, query, chatID, userID)
	return &participant, err
}

func (r *ParticipantRepository) Update(ctx context.Context, participant *domain.Participant) error {
	query := `UPDATE participants SET role = $1, last_read_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, participant.Role, participant.LastReadAt, participant.ID)
	return err
}

func (r *ParticipantRepository) UpdateLastRead(ctx context.Context, chatID, userID uuid.UUID, lastReadAt time.Time) error {
	query := `UPDATE participants SET last_read_at = $1 WHERE chat_id = $2 AND user_id = $3`
	_, err := r.db.ExecContext(ctx, query, lastReadAt, chatID, userID)
	return err
}

func (r *ParticipantRepository) GetUnreadCount(ctx context.Context, chatID, userID uuid.UUID, lastReadAt time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE chat_id = $1 
		AND user_id != $2 
		AND deleted_at IS NULL
		AND created_at > $3
	`
	var count int
	err := r.db.GetContext(ctx, &count, query, chatID, userID, lastReadAt)
	return count, err
}

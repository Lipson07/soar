package postgres

import (
	"context"
	"database/sql"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CallRepository struct {
	db *sqlx.DB
}

func NewCallRepository(db *sqlx.DB) repository.CallRepository {
	return &CallRepository{db: db}
}

func (r *CallRepository) Create(ctx context.Context, call *domain.Call) error {
	query := `
		INSERT INTO calls (id, chat_id, caller_id, callee_id, type, status, room_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		call.ID, call.ChatID, call.CallerID, call.CalleeID,
		call.Type, call.Status, call.RoomID, call.CreatedAt,
	)
	return err
}

func (r *CallRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Call, error) {
	query := `SELECT id, chat_id, caller_id, callee_id, type, status, room_id, started_at, ended_at, created_at FROM calls WHERE id = $1`
	var call domain.Call
	err := r.db.GetContext(ctx, &call, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &call, err
}

func (r *CallRepository) GetByRoomID(ctx context.Context, roomID string) (*domain.Call, error) {
	query := `SELECT id, chat_id, caller_id, callee_id, type, status, room_id, started_at, ended_at, created_at FROM calls WHERE room_id = $1`
	var call domain.Call
	err := r.db.GetContext(ctx, &call, query, roomID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &call, err
}

func (r *CallRepository) GetActiveCall(ctx context.Context, chatID uuid.UUID) (*domain.Call, error) {
	query := `SELECT id, chat_id, caller_id, callee_id, type, status, room_id, started_at, ended_at, created_at FROM calls WHERE chat_id = $1 AND status IN ('pending', 'active') LIMIT 1`
	var call domain.Call
	err := r.db.GetContext(ctx, &call, query, chatID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &call, err
}

func (r *CallRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CallStatus) error {
	query := `UPDATE calls SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *CallRepository) UpdateEnded(ctx context.Context, id uuid.UUID, endedAt time.Time) error {
	query := `UPDATE calls SET status = 'ended', ended_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, endedAt, id)
	return err
}

func (r *CallRepository) GetUserCalls(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Call, error) {
	query := `SELECT id, chat_id, caller_id, callee_id, type, status, room_id, started_at, ended_at, created_at FROM calls WHERE caller_id = $1 OR callee_id = $1 ORDER BY created_at DESC LIMIT $2`
	var calls []*domain.Call
	err := r.db.SelectContext(ctx, &calls, query, userID, limit)
	return calls, err
}

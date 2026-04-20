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

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) repository.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *domain.UserSession) error {
	query := `
		INSERT INTO user_sessions (
			user_id, session_token, device_info, device_type,
			ip_address, location, user_agent, last_active, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		session.UserID,
		session.SessionToken,
		session.DeviceInfo,
		session.DeviceType,
		session.IPAddress,
		session.Location,
		session.UserAgent,
		time.Now(),
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *SessionRepository) GetUserActiveSessions(ctx context.Context, userID uuid.UUID) ([]domain.UserSession, error) {
	query := `
		SELECT id, user_id, session_token, device_info, device_type,
		       ip_address, location, user_agent, last_active, created_at, is_active
		FROM user_sessions
		WHERE user_id = $1 AND is_active = true
		ORDER BY last_active DESC
	`

	var sessions []domain.UserSession
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	return sessions, err
}

func (r *SessionRepository) TerminateSession(ctx context.Context, sessionID int64, userID uuid.UUID) error {
	query := `
		UPDATE user_sessions
		SET is_active = false
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, sessionID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SessionRepository) TerminateAllOtherSessions(ctx context.Context, userID uuid.UUID, currentSessionToken string) error {
	query := `
		UPDATE user_sessions
		SET is_active = false
		WHERE user_id = $1 AND session_token != $2 AND is_active = true
	`

	_, err := r.db.ExecContext(ctx, query, userID, currentSessionToken)
	return err
}

func (r *SessionRepository) UpdateSessionActivity(ctx context.Context, sessionToken string) error {
	query := `
		UPDATE user_sessions
		SET last_active = $1
		WHERE session_token = $2 AND is_active = true
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), sessionToken)
	return err
}

func (r *SessionRepository) GetSessionByToken(ctx context.Context, token string) (*domain.UserSession, error) {
	query := `
		SELECT id, user_id, session_token, device_info, device_type,
		       ip_address, location, user_agent, last_active, created_at, is_active
		FROM user_sessions
		WHERE session_token = $1 AND is_active = true
	`

	var session domain.UserSession
	err := r.db.GetContext(ctx, &session, query, token)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

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

type SecurityRepository struct {
	db *sqlx.DB
}

func NewSecurityRepository(db *sqlx.DB) repository.SecurityRepository {
	return &SecurityRepository{db: db}
}

func (r *SecurityRepository) GetUserSecuritySettings(ctx context.Context, userID uuid.UUID) (*domain.SecuritySettings, error) {
	query := `
		SELECT id, user_id, two_factor_enabled, two_factor_secret, 
		       biometric_enabled, end_to_end_encryption, screen_security,
		       login_alerts, created_at, updated_at
		FROM security_settings
		WHERE user_id = $1
	`

	var settings domain.SecuritySettings
	err := r.db.GetContext(ctx, &settings, query, userID)
	if err == sql.ErrNoRows {
		return r.CreateDefaultSettings(ctx, userID)
	}
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (r *SecurityRepository) CreateDefaultSettings(ctx context.Context, userID uuid.UUID) (*domain.SecuritySettings, error) {
	query := `
		INSERT INTO security_settings (
			user_id, two_factor_enabled, biometric_enabled,
			end_to_end_encryption, screen_security, login_alerts
		) VALUES ($1, false, false, true, true, true)
		ON CONFLICT (user_id) DO UPDATE SET
			user_id = EXCLUDED.user_id
		RETURNING id, user_id, two_factor_enabled, two_factor_secret,
		          biometric_enabled, end_to_end_encryption, screen_security,
		          login_alerts, created_at, updated_at
	`

	var settings domain.SecuritySettings
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.TwoFactorEnabled,
		&settings.TwoFactorSecret,
		&settings.BiometricEnabled,
		&settings.EndToEndEncryption,
		&settings.ScreenSecurity,
		&settings.LoginAlerts,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	return &settings, err
}

func (r *SecurityRepository) UpdateSecuritySettings(ctx context.Context, settings *domain.SecuritySettings) error {
	query := `
		INSERT INTO security_settings (
			user_id, two_factor_enabled, two_factor_secret,
			biometric_enabled, end_to_end_encryption,
			screen_security, login_alerts, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) DO UPDATE SET
			two_factor_enabled = EXCLUDED.two_factor_enabled,
			two_factor_secret = EXCLUDED.two_factor_secret,
			biometric_enabled = EXCLUDED.biometric_enabled,
			end_to_end_encryption = EXCLUDED.end_to_end_encryption,
			screen_security = EXCLUDED.screen_security,
			login_alerts = EXCLUDED.login_alerts,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		settings.UserID,
		settings.TwoFactorEnabled,
		settings.TwoFactorSecret,
		settings.BiometricEnabled,
		settings.EndToEndEncryption,
		settings.ScreenSecurity,
		settings.LoginAlerts,
		time.Now(),
	)

	return err
}

func (r *SecurityRepository) EnableTwoFactor(ctx context.Context, userID uuid.UUID, secret string) error {
	query := `
		INSERT INTO security_settings (
			user_id, two_factor_enabled, two_factor_secret, updated_at
		) VALUES ($1, true, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			two_factor_enabled = EXCLUDED.two_factor_enabled,
			two_factor_secret = EXCLUDED.two_factor_secret,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, userID, secret, time.Now())
	return err
}

func (r *SecurityRepository) DisableTwoFactor(ctx context.Context, userID uuid.UUID) error {
	query := `
		INSERT INTO security_settings (
			user_id, two_factor_enabled, two_factor_secret, updated_at
		) VALUES ($1, false, '', $2)
		ON CONFLICT (user_id) DO UPDATE SET
			two_factor_enabled = EXCLUDED.two_factor_enabled,
			two_factor_secret = EXCLUDED.two_factor_secret,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	return err
}

func (r *SecurityRepository) CreateAuditLog(ctx context.Context, log *domain.SecurityAuditLog) error {
	query := `
		INSERT INTO security_audit_logs (user_id, action, ip_address, user_agent, details)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.UserID,
		log.Action,
		log.IPAddress,
		log.UserAgent,
		log.Details,
	)

	return err
}

func (r *SecurityRepository) GetUserAuditLogs(ctx context.Context, userID uuid.UUID, limit int) ([]domain.SecurityAuditLog, error) {
	query := `
		SELECT id, user_id, action, ip_address, user_agent, details, created_at
		FROM security_audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	var logs []domain.SecurityAuditLog
	err := r.db.SelectContext(ctx, &logs, query, userID, limit)
	return logs, err
}

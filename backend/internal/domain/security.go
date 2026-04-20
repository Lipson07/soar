package domain

import (
	"time"

	"github.com/google/uuid"
)

type SecuritySettings struct {
	ID                 int64     `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	TwoFactorEnabled   bool      `json:"two_factor_enabled" db:"two_factor_enabled"`
	TwoFactorSecret    string    `json:"-" db:"two_factor_secret"`
	BiometricEnabled   bool      `json:"biometric_enabled" db:"biometric_enabled"`
	EndToEndEncryption bool      `json:"end_to_end_encryption" db:"end_to_end_encryption"`
	ScreenSecurity     bool      `json:"screen_security" db:"screen_security"`
	LoginAlerts        bool      `json:"login_alerts" db:"login_alerts"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type UserSession struct {
	ID           int64     `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	SessionToken string    `json:"session_token" db:"session_token"`
	DeviceInfo   string    `json:"device_info" db:"device_info"`
	DeviceType   string    `json:"device_type" db:"device_type"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	Location     string    `json:"location" db:"location"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	LastActive   time.Time `json:"last_active" db:"last_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsCurrent    bool      `json:"is_current" db:"-"`
}

type SecurityAuditLog struct {
	ID        int64     `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Details   string    `json:"details" db:"details"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type SecurityReport struct {
	Settings      *SecuritySettings  `json:"settings"`
	Sessions      []UserSession      `json:"sessions"`
	AuditLogs     []SecurityAuditLog `json:"audit_logs,omitempty"`
	GeneratedAt   time.Time          `json:"generated_at"`
	SecurityScore int                `json:"security_score"`
}

type TwoFactorSetup struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

type TwoFactorVerifyRequest struct {
	Code   string `json:"code" validate:"required,len=6"`
	Secret string `json:"secret" validate:"required"`
}

type UpdateSecuritySettingsRequest struct {
	TwoFactorEnabled   *bool `json:"two_factor_enabled,omitempty"`
	BiometricEnabled   *bool `json:"biometric_enabled,omitempty"`
	EndToEndEncryption *bool `json:"end_to_end_encryption,omitempty"`
	ScreenSecurity     *bool `json:"screen_security,omitempty"`
	LoginAlerts        *bool `json:"login_alerts,omitempty"`
}

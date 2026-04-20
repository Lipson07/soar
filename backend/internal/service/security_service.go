package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type securityServiceImpl struct {
	securityRepo repository.SecurityRepository
	sessionRepo  repository.SessionRepository
}

func NewSecurityService(
	securityRepo repository.SecurityRepository,
	sessionRepo repository.SessionRepository,
) SecurityService {
	return &securityServiceImpl{
		securityRepo: securityRepo,
		sessionRepo:  sessionRepo,
	}
}

func (s *securityServiceImpl) GetUserSettings(ctx context.Context, userID uuid.UUID) (*domain.SecuritySettings, error) {
	return s.securityRepo.GetUserSecuritySettings(ctx, userID)
}

func (s *securityServiceImpl) UpdateSettings(ctx context.Context, userID uuid.UUID, req *domain.UpdateSecuritySettingsRequest) error {
	settings, err := s.securityRepo.GetUserSecuritySettings(ctx, userID)
	if err != nil {
		return err
	}

	if req.TwoFactorEnabled != nil {
		settings.TwoFactorEnabled = *req.TwoFactorEnabled
	}
	if req.BiometricEnabled != nil {
		settings.BiometricEnabled = *req.BiometricEnabled
	}
	if req.EndToEndEncryption != nil {
		settings.EndToEndEncryption = *req.EndToEndEncryption
	}
	if req.ScreenSecurity != nil {
		settings.ScreenSecurity = *req.ScreenSecurity
	}
	if req.LoginAlerts != nil {
		settings.LoginAlerts = *req.LoginAlerts
	}

	err = s.securityRepo.UpdateSecuritySettings(ctx, settings)
	if err != nil {
		return err
	}

	s.createAuditLog(ctx, userID, "settings_updated", "Security settings updated")
	return nil
}

func (s *securityServiceImpl) SetupTwoFactor(ctx context.Context, userID uuid.UUID, username string) (*domain.TwoFactorSetup, error) {
	secret := generateSecret()

	issuer := "Messenger"
	otpAuthURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, username, secret, issuer)

	backupCodes := generateBackupCodes()

	return &domain.TwoFactorSetup{
		Secret:      secret,
		QRCode:      otpAuthURL,
		BackupCodes: backupCodes,
	}, nil
}

func (s *securityServiceImpl) VerifyAndEnableTwoFactor(ctx context.Context, userID uuid.UUID, code string, secret string) (bool, error) {
	valid := validateTOTP(code, secret)
	if !valid {
		return false, nil
	}

	err := s.securityRepo.EnableTwoFactor(ctx, userID, secret)
	if err != nil {
		return false, err
	}

	s.createAuditLog(ctx, userID, "2fa_enabled", "Two-factor authentication enabled")
	return true, nil
}

func (s *securityServiceImpl) DisableTwoFactor(ctx context.Context, userID uuid.UUID) error {
	err := s.securityRepo.DisableTwoFactor(ctx, userID)
	if err != nil {
		return err
	}

	s.createAuditLog(ctx, userID, "2fa_disabled", "Two-factor authentication disabled")
	return nil
}

func (s *securityServiceImpl) VerifyTwoFactorCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	settings, err := s.securityRepo.GetUserSecuritySettings(ctx, userID)
	if err != nil {
		return false, err
	}

	if !settings.TwoFactorEnabled {
		return true, nil
	}

	valid := validateTOTP(code, settings.TwoFactorSecret)

	action := "2fa_login_success"
	if !valid {
		action = "2fa_login_failed"
	}
	s.createAuditLog(ctx, userID, action, fmt.Sprintf("2FA verification: %v", valid))

	return valid, nil
}

func (s *securityServiceImpl) GetUserSessions(ctx context.Context, userID uuid.UUID, currentToken string) ([]domain.UserSession, error) {
	sessions, err := s.sessionRepo.GetUserActiveSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range sessions {
		if sessions[i].SessionToken == currentToken {
			sessions[i].IsCurrent = true
		}
	}

	return sessions, nil
}

func (s *securityServiceImpl) TerminateSession(ctx context.Context, userID uuid.UUID, sessionID int64, currentToken string) error {
	sessions, err := s.sessionRepo.GetUserActiveSessions(ctx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.ID == sessionID && session.SessionToken == currentToken {
			return fmt.Errorf("cannot terminate current session")
		}
	}

	err = s.sessionRepo.TerminateSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}

	s.createAuditLog(ctx, userID, "session_terminated", fmt.Sprintf("Session %d terminated", sessionID))
	return nil
}

func (s *securityServiceImpl) TerminateAllOtherSessions(ctx context.Context, userID uuid.UUID, currentToken string) error {
	err := s.sessionRepo.TerminateAllOtherSessions(ctx, userID, currentToken)
	if err != nil {
		return err
	}

	s.createAuditLog(ctx, userID, "all_sessions_terminated", "All other sessions terminated")
	return nil
}

func (s *securityServiceImpl) GenerateSecurityReport(ctx context.Context, userID uuid.UUID) (*domain.SecurityReport, error) {
	settings, err := s.securityRepo.GetUserSecuritySettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessions, err := s.sessionRepo.GetUserActiveSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	auditLogs, err := s.securityRepo.GetUserAuditLogs(ctx, userID, 50)
	if err != nil {
		return nil, err
	}

	securityScore := calculateSecurityScore(settings)

	return &domain.SecurityReport{
		Settings:      settings,
		Sessions:      sessions,
		AuditLogs:     auditLogs,
		GeneratedAt:   time.Now(),
		SecurityScore: securityScore,
	}, nil
}

func (s *securityServiceImpl) CreateSession(ctx context.Context, userID uuid.UUID, deviceInfo, deviceType, ipAddress, userAgent string) (*domain.UserSession, error) {
	session := &domain.UserSession{
		UserID:       userID,
		SessionToken: generateSessionToken(),
		DeviceInfo:   deviceInfo,
		DeviceType:   deviceType,
		IPAddress:    ipAddress,
		Location:     getLocationFromIP(ipAddress),
		UserAgent:    userAgent,
	}

	err := s.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	settings, _ := s.securityRepo.GetUserSecuritySettings(ctx, userID)
	if settings != nil && settings.LoginAlerts {
		s.sendLoginAlert(ctx, userID, session)
	}

	s.createAuditLog(ctx, userID, "login", fmt.Sprintf("New session from %s", deviceInfo))

	return session, nil
}

func (s *securityServiceImpl) createAuditLog(ctx context.Context, userID uuid.UUID, action, details string) {
	log := &domain.SecurityAuditLog{
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}

	if ip, ok := ctx.Value("ip_address").(string); ok {
		log.IPAddress = ip
	}
	if ua, ok := ctx.Value("user_agent").(string); ok {
		log.UserAgent = ua
	}

	_ = s.securityRepo.CreateAuditLog(ctx, log)
}

func (s *securityServiceImpl) sendLoginAlert(ctx context.Context, userID uuid.UUID, session *domain.UserSession) {
}

func validateTOTP(code string, secret string) bool {
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return false
	}

	counter := uint64(time.Now().Unix() / 30)

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	h := hmac.New(sha1.New, secretBytes)
	h.Write(buf)
	hash := h.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	otp := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF % 1000000

	return code == fmt.Sprintf("%06d", otp)
}

func generateSecret() string {
	secret := make([]byte, 20)
	rand.Read(secret)
	return base32.StdEncoding.EncodeToString(secret)
}

func generateBackupCodes() []string {
	codes := make([]string, 8)
	for i := 0; i < 8; i++ {
		b := make([]byte, 6)
		rand.Read(b)
		codes[i] = fmt.Sprintf("%x", b)[:8]
	}
	return codes
}

func generateSessionToken() string {
	return uuid.New().String()
}

func calculateSecurityScore(settings *domain.SecuritySettings) int {
	score := 0
	if settings.TwoFactorEnabled {
		score += 30
	}
	if settings.BiometricEnabled {
		score += 20
	}
	if settings.EndToEndEncryption {
		score += 25
	}
	if settings.ScreenSecurity {
		score += 15
	}
	if settings.LoginAlerts {
		score += 10
	}
	return score
}

func getLocationFromIP(ip string) string {
	return "Unknown"
}

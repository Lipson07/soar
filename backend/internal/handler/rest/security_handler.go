package rest

import (
	"net/http"
	"strconv"

	"backend/internal/domain"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecurityHandler struct {
	securityService service.SecurityService
}

func NewSecurityHandler(securityService service.SecurityService) *SecurityHandler {
	return &SecurityHandler{
		securityService: securityService,
	}
}

func (h *SecurityHandler) getUserID(c *gin.Context) uuid.UUID {
	val, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}
	switch v := val.(type) {
	case uuid.UUID:
		return v
	case string:
		id, _ := uuid.Parse(v)
		return id
	default:
		return uuid.Nil
	}
}

func (h *SecurityHandler) getString(c *gin.Context, key string) string {
	val, exists := c.Get(key)
	if !exists {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

func (h *SecurityHandler) GetSettings(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	settings, err := h.securityService.GetUserSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *SecurityHandler) UpdateSettings(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req domain.UpdateSecuritySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.securityService.UpdateSettings(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *SecurityHandler) SetupTwoFactor(c *gin.Context) {
	userID := h.getUserID(c)
	username := h.getString(c, "username")

	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	setup, err := h.securityService.SetupTwoFactor(c.Request.Context(), userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setup)
}

func (h *SecurityHandler) VerifyTwoFactor(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req domain.TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := h.securityService.VerifyAndEnableTwoFactor(c.Request.Context(), userID, req.Code, req.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enabled": true})
}

func (h *SecurityHandler) DisableTwoFactor(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := h.securityService.DisableTwoFactor(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enabled": false})
}

func (h *SecurityHandler) GetSessions(c *gin.Context) {
	userID := h.getUserID(c)
	currentToken := h.getString(c, "session_token")

	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	sessions, err := h.securityService.GetUserSessions(c.Request.Context(), userID, currentToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *SecurityHandler) TerminateSession(c *gin.Context) {
	userID := h.getUserID(c)
	currentToken := h.getString(c, "session_token")

	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := h.securityService.TerminateSession(c.Request.Context(), userID, sessionID, currentToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "terminated"})
}

func (h *SecurityHandler) TerminateAllOtherSessions(c *gin.Context) {
	userID := h.getUserID(c)
	currentToken := h.getString(c, "session_token")

	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := h.securityService.TerminateAllOtherSessions(c.Request.Context(), userID, currentToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "terminated"})
}

func (h *SecurityHandler) GetSecurityReport(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	report, err := h.securityService.GenerateSecurityReport(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, report)
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=security-report.json")
	c.JSON(http.StatusOK, report)
}

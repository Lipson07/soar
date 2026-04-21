package rest

import (
	"net/http"

	"backend/internal/domain"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CallHandler struct {
	callService service.CallService
}

func NewCallHandler(callService service.CallService) *CallHandler {
	return &CallHandler{
		callService: callService,
	}
}

func (h *CallHandler) StartCall(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		ChatID   uuid.UUID `json:"chat_id"`
		CalleeID uuid.UUID `json:"callee_id"`
		Type     string    `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	callType := domain.CallType(req.Type)
	if callType != domain.CallTypeAudio && callType != domain.CallTypeVideo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call type"})
		return
	}

	call, err := h.callService.StartCall(c.Request.Context(), req.ChatID, userID, req.CalleeID, callType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, call)
}

func (h *CallHandler) AcceptCall(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("callId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	call, err := h.callService.AcceptCall(c.Request.Context(), callID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, call)
}

func (h *CallHandler) RejectCall(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("callId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	call, err := h.callService.RejectCall(c.Request.Context(), callID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, call)
}

func (h *CallHandler) EndCall(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("callId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	call, err := h.callService.EndCall(c.Request.Context(), callID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, call)
}

func (h *CallHandler) GetActiveCall(c *gin.Context) {
	chatID, err := uuid.Parse(c.Query("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	call, err := h.callService.GetActiveCall(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, call)
}

func (h *CallHandler) GetUserCalls(c *gin.Context) {
	userID := getUserIDFromContext(c)
	limit := 50

	calls, err := h.callService.GetUserCalls(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, calls)
}

func getUserIDFromContext(c *gin.Context) uuid.UUID {
	val, _ := c.Get("user_id")
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

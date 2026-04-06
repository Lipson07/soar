package rest

import (
	"backend/internal/domain"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// SendMessage отправляет сообщение
// @Summary Отправить сообщение
// @Description Отправляет сообщение в чат
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chat_id query string true "ID чата"
// @Param request body domain.SendMessageRequest true "Текст сообщения"
// @Success 201 {object} domain.Message
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	var req domain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Исправляем проблему с конвертацией userID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	// Конвертируем userID из string в uuid.UUID
	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "некорректный формат ID пользователя"})
			return
		}
	case uuid.UUID:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неподдерживаемый тип ID пользователя"})
		return
	}

	message, err := h.messageService.SendMessage(c.Request.Context(), chatID, &req, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// GetMessages возвращает сообщения чата
// @Summary Получить сообщения
// @Description Возвращает сообщения чата с пагинацией
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chat_id query string true "ID чата"
// @Param limit query int false "Лимит (default 50)"
// @Param offset query int false "Смещение (default 0)"
// @Success 200 {array} domain.Message
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Исправляем проблему с конвертацией userID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "некорректный формат ID пользователя"})
			return
		}
	case uuid.UUID:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неподдерживаемый тип ID пользователя"})
		return
	}

	messages, err := h.messageService.GetMessages(c.Request.Context(), chatID, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// EditMessage редактирует сообщение
// @Summary Редактировать сообщение
// @Description Редактирует текст отправленного сообщения
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message_id query string true "ID сообщения"
// @Param request body object true "Новый текст"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/messages [put]
func (h *MessageHandler) EditMessage(c *gin.Context) {
	messageIDStr := c.Query("message_id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID сообщения"})
		return
	}

	var req struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Исправляем проблему с конвертацией userID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "некорректный формат ID пользователя"})
			return
		}
	case uuid.UUID:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неподдерживаемый тип ID пользователя"})
		return
	}

	err = h.messageService.EditMessage(c.Request.Context(), messageID, req.Text, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "сообщение отредактировано"})
}

// DeleteMessage удаляет сообщение
// @Summary Удалить сообщение
// @Description Удаляет сообщение (своё или если вы админ)
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message_id query string true "ID сообщения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/messages [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	messageIDStr := c.Query("message_id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID сообщения"})
		return
	}

	// Исправляем проблему с конвертацией userID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	var userID uuid.UUID
	switch v := userIDInterface.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "некорректный формат ID пользователя"})
			return
		}
	case uuid.UUID:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неподдерживаемый тип ID пользователя"})
		return
	}

	err = h.messageService.DeleteMessage(c.Request.Context(), messageID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "сообщение удалено"})
}

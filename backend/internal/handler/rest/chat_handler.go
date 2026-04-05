package rest

import (
	"backend/internal/domain"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// CreatePrivateChat создает личный чат
// @Summary      Создать личный чат
// @Description  Создает приватный чат между текущим пользователем и указанным пользователем
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreatePrivateChatRequest true "ID собеседника"
// @Success      201 {object} domain.Chat
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/chats/private [post]
func (h *ChatHandler) CreatePrivateChat(c *gin.Context) {
	var req domain.CreatePrivateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	// Конвертируем строку в UUID (так как middleware сохраняет строку)
	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	creatorID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный формат ID пользователя"})
		return
	}

	chat, err := h.chatService.CreatePrivateChat(c.Request.Context(), creatorID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

// CreateGroupChat создает групповой чат
// @Summary      Создать групповой чат
// @Description  Создает групповой чат с указанными участниками
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateGroupChatRequest true "Данные группы"
// @Success      201 {object} domain.Chat
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/chats/group [post]
func (h *ChatHandler) CreateGroupChat(c *gin.Context) {
	var req domain.CreateGroupChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	creatorID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный формат ID пользователя"})
		return
	}

	chat, err := h.chatService.CreateGroupChat(c.Request.Context(), &req, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

// GetUserChats возвращает все чаты пользователя
// @Summary      Получить чаты пользователя
// @Description  Возвращает все чаты, в которых участвует текущий пользователь
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} domain.ChatResponse
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/chats [get]
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный формат ID пользователя"})
		return
	}

	chats, err := h.chatService.GetUserChats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if chats == nil {
		c.JSON(http.StatusOK, []domain.ChatResponse{})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// GetChatByID возвращает чат по ID
// @Summary      Получить чат по ID
// @Description  Возвращает информацию о чате по его ID
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID чата"
// @Success      200 {object} domain.Chat
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/chats/{id} [get]
func (h *ChatHandler) GetChatByID(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	chat, err := h.chatService.GetChatByID(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// UpdateChat обновляет чат
// @Summary      Обновить чат
// @Description  Обновляет данные чата по ID (только для администраторов)
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID чата"
// @Param        request body domain.UpdateChatRequest true "Данные для обновления"
// @Success      200 {object} domain.Chat
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/chats/{id} [put]
func (h *ChatHandler) UpdateChat(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	var req domain.UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := h.chatService.UpdateChat(c.Request.Context(), chatID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// DeleteChat удаляет чат
// @Summary      Удалить чат
// @Description  Удаляет чат по ID (только для администраторов)
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID чата"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/chats/{id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	err = h.chatService.DeleteChat(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "чат успешно удален"})
}

// GetAllChats возвращает все чаты
// @Summary      Получить все чаты
// @Description  Возвращает список всех чатов (только для администраторов)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} domain.Chat
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/admin/chats [get]
func (h *ChatHandler) GetAllChats(c *gin.Context) {
	chats, err := h.chatService.GetAllChats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}
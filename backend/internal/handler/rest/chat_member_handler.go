package rest

import (
	"net/http"
	"strconv"

	"myapp/internal/domain"
	"myapp/internal/service"

	"github.com/gin-gonic/gin"
)

type ChatMemberHandler struct {
	chatMemberService service.ChatMemberService
}

func NewChatMemberHandler(chatMemberService service.ChatMemberService) *ChatMemberHandler {
	return &ChatMemberHandler{
		chatMemberService: chatMemberService,
	}
}

// AddMember godoc
// @Summary      Добавить участника в чат
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Param        request body domain.AddMemberRequest true "Данные участника"
// @Success      201 {object} domain.ChatMember
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members [post]
func (h *ChatMemberHandler) AddMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	var req domain.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetInt64("user_id")

	member, err := h.chatMemberService.AddMember(c, chatID, &req, currentUserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// AddMembers godoc
// @Summary      Массовое добавление участников
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Param        request body object true "Список ID пользователей" example({"user_ids":[1,2,3]})
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members/bulk [post]
func (h *ChatMemberHandler) AddMembers(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	var req struct {
		UserIDs []int64 `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetInt64("user_id")

	if err := h.chatMemberService.AddMembers(c, chatID, req.UserIDs, currentUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "участники успешно добавлены"})
}

// GetChatMembers godoc
// @Summary      Получить всех участников чата
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Success      200 {array} domain.ChatMember
// @Failure      400 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members [get]
func (h *ChatMemberHandler) GetChatMembers(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	members, err := h.chatMemberService.GetChatMembers(c, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetMember godoc
// @Summary      Получить информацию об участнике
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Param        user_id path int true "ID пользователя"
// @Success      200 {object} domain.ChatMember
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members/{user_id} [get]
func (h *ChatMemberHandler) GetMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID пользователя"})
		return
	}

	member, err := h.chatMemberService.GetMember(c, chatID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

// UpdateMemberRole godoc
// @Summary      Обновить роль участника
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Param        user_id path int true "ID пользователя"
// @Param        request body domain.UpdateMemberRoleRequest true "Новая роль"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members/{user_id}/role [put]
func (h *ChatMemberHandler) UpdateMemberRole(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID пользователя"})
		return
	}

	var req domain.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := c.GetInt64("user_id")

	if err := h.chatMemberService.UpdateMemberRole(c, chatID, userID, &req, currentUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "роль успешно обновлена"})
}

// UpdateLastRead godoc
// @Summary      Обновить время последнего прочтения
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /chats/{chat_id}/read [post]
func (h *ChatMemberHandler) UpdateLastRead(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	currentUserID := c.GetInt64("user_id")

	if err := h.chatMemberService.UpdateLastRead(c, chatID, currentUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "время прочтения обновлено"})
}

// RemoveMember godoc
// @Summary      Удалить участника из чата
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Param        user_id path int true "ID пользователя"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members/{user_id} [delete]
func (h *ChatMemberHandler) RemoveMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID пользователя"})
		return
	}

	currentUserID := c.GetInt64("user_id")

	if err := h.chatMemberService.RemoveMember(c, chatID, userID, currentUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "участник успешно удален"})
}

// LeaveChat godoc
// @Summary      Выйти из чата
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /chats/{chat_id}/leave [post]
func (h *ChatMemberHandler) LeaveChat(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	currentUserID := c.GetInt64("user_id")

	if err := h.chatMemberService.LeaveChat(c, chatID, currentUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "вы вышли из чата"})
}

// KickMember godoc
// @Summary      Исключить участника из чата
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Param        user_id path int true "ID пользователя"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members/{user_id}/kick [post]
func (h *ChatMemberHandler) KickMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID пользователя"})
		return
	}

	currentUserID := c.GetInt64("user_id")

	if err := h.chatMemberService.KickMember(c, chatID, userID, currentUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "участник исключен из чата"})
}

// GetUserChats godoc
// @Summary      Получить все чаты пользователя
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Success      200 {array} domain.Chat
// @Failure      500 {object} map[string]interface{}
// @Router       /users/chats [get]
func (h *ChatMemberHandler) GetUserChats(c *gin.Context) {
	currentUserID := c.GetInt64("user_id")
	if currentUserID == 0 {
		currentUserID = 1
	}

	chats, err := h.chatMemberService.GetUserChats(c, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if chats == nil {
		c.JSON(http.StatusOK, []domain.Chat{})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// GetMemberCount godoc
// @Summary      Получить количество участников чата
// @Tags         участники-чатов
// @Accept       json
// @Produce      json
// @Param        chat_id path int true "ID чата"
// @Success      200 {object} map[string]int
// @Failure      400 {object} map[string]interface{}
// @Router       /chats/{chat_id}/members/count [get]
func (h *ChatMemberHandler) GetMemberCount(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	count, err := h.chatMemberService.GetMemberCount(c, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

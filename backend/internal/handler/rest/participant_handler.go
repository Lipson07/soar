package rest

import (
	"backend/internal/domain"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ParticipantHandler struct {
	participantService service.ParticipantService
}

func NewParticipantHandler(participantService service.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{
		participantService: participantService,
	}
}

// AddParticipants добавляет участников в чат
// @Summary      Добавить участников
// @Description  Добавляет пользователей в групповой чат (только для администраторов)
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Param        request body domain.AddParticipantsRequest true "Список ID пользователей"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /api/participants [post]
func (h *ParticipantHandler) AddParticipants(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	var req domain.AddParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	adderID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	err = h.participantService.AddParticipants(c.Request.Context(), chatID, req.UserIDs, adderID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "участники успешно добавлены"})
}

// RemoveParticipant удаляет участника из чата
// @Summary      Удалить участника
// @Description  Удаляет пользователя из группового чата (только для администраторов)
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Param        user_id query string true "ID пользователя"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /api/participants [delete]
func (h *ParticipantHandler) RemoveParticipant(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID пользователя"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	removerID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	err = h.participantService.RemoveParticipant(c.Request.Context(), chatID, userID, removerID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "участник успешно удален"})
}

// LeaveChat выход из чата
// @Summary      Выйти из чата
// @Description  Пользователь выходит из группового чата
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /api/participants/leave [post]
func (h *ParticipantHandler) LeaveChat(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	err = h.participantService.LeaveChat(c.Request.Context(), chatID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "вы вышли из чата"})
}

// GetChatParticipants возвращает участников чата
// @Summary      Получить участников чата
// @Description  Возвращает список всех участников чата
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Success      200 {array} domain.Participant
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /api/participants [get]
func (h *ParticipantHandler) GetChatParticipants(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	participants, err := h.participantService.GetChatParticipants(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participants)
}

// UpdateRole обновляет роль участника
// @Summary      Обновить роль
// @Description  Обновляет роль участника в чате (admin/member)
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Param        request body domain.UpdateRoleRequest true "Данные для обновления роли"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /api/participants/role [put]
func (h *ParticipantHandler) UpdateRole(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	var req domain.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	updaterID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	err = h.participantService.UpdateRole(c.Request.Context(), chatID, req.UserID, updaterID, req.Role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "роль успешно обновлена"})
}

// UpdateLastRead обновляет время последнего прочтения
// @Summary      Обновить прочтение
// @Description  Обновляет время последнего прочтения сообщений в чате
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /api/participants/read [put]
func (h *ParticipantHandler) UpdateLastRead(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	err = h.participantService.UpdateLastRead(c.Request.Context(), chatID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "время прочтения обновлено"})
}

// GetUnreadCount возвращает количество непрочитанных сообщений
// @Summary      Непрочитанные сообщения
// @Description  Возвращает количество непрочитанных сообщений в чате
// @Tags         participants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id query string true "ID чата"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /api/participants/unread [get]
func (h *ParticipantHandler) GetUnreadCount(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип ID пользователя"})
		return
	}

	count, err := h.participantService.GetUnreadCount(c.Request.Context(), chatID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
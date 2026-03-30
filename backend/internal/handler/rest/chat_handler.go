package rest

import (
	"net/http"
	"strconv"

	"myapp/internal/domain"
	"myapp/internal/service"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// Create создает новый чат
// @Summary      Создание чата
// @Description  Создает новый чат (private, group или channel)
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        request body domain.CreateChatRequest true "Данные чата"
// @Success      201  {object}  domain.Chat
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /chats [post]
func (h *ChatHandler) Create(c *gin.Context) {
	var req domain.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := h.chatService.Create(c, &req)
	if err != nil {
		if err == domain.ErrChatNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == domain.ErrInvalidChatType {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

// GetByID возвращает чат по ID
// @Summary      Получить чат
// @Description  Возвращает чат по его ID
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID чата"
// @Success      200  {object}  domain.Chat
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /chats/{id} [get]
func (h *ChatHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	chat, err := h.chatService.GetByID(c, id)
	if err != nil {
		if err == domain.ErrChatNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// GetByName возвращает чат по имени
// @Summary      Получить чат по имени
// @Description  Возвращает чат по его названию
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        name   query      string  true  "Название чата"
// @Success      200  {object}  domain.Chat
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /chats/by-name [get]
func (h *ChatHandler) GetByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не указано название чата"})
		return
	}

	chat, err := h.chatService.GetByName(c, name)
	if err != nil {
		if err == domain.ErrChatNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// GetAll возвращает все чаты
// @Summary      Получить все чаты
// @Description  Возвращает список всех чатов
// @Tags         chats
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Chat
// @Failure      500  {object}  map[string]interface{}
// @Router       /chats [get]
func (h *ChatHandler) GetAll(c *gin.Context) {
	chats, err := h.chatService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// Update обновляет чат
// @Summary      Обновить чат
// @Description  Обновляет данные существующего чата
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id      path      int                       true  "ID чата"
// @Param        request body      domain.UpdateChatRequest true  "Новые данные"
// @Success      200     {object}  domain.Chat
// @Failure      400     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Failure      409     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /chats/{id} [put]
func (h *ChatHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	var req domain.UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := h.chatService.Update(c, id, &req)
	if err != nil {
		if err == domain.ErrChatNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == domain.ErrChatNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// Delete удаляет чат
// @Summary      Удалить чат
// @Description  Удаляет чат по ID
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID чата"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /chats/{id} [delete]
func (h *ChatHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.chatService.Delete(c, id); err != nil {
		if err == domain.ErrChatNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "чат удален"})
}
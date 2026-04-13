package rest

import (
	"backend/internal/domain"
	"backend/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func (h *MessageHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл не предоставлен"})
		return
	}

	maxSize := int64(10 << 20)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл слишком большой (максимум 10MB)"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ext

	workDir, _ := os.Getwd()
	mimeType := file.Header.Get("Content-Type")

	var uploadDir string
	var fileURL string

	if strings.HasPrefix(mimeType, "image/") {
		uploadDir = filepath.Join(workDir, "uploads", "images")
		fileURL = "/uploads/images/" + filename
	} else {
		uploadDir = filepath.Join(workDir, "uploads", "files")
		fileURL = "/uploads/files/" + filename
	}

	os.MkdirAll(uploadDir, 0755)

	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения файла"})
		return
	}

	response := domain.UploadFileResponse{
		FileURL:  fileURL,
		FileName: file.Filename,
		FileSize: file.Size,
		MimeType: mimeType,
	}

	c.JSON(http.StatusOK, response)
}

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

	message, err := h.messageService.SendMessage(c.Request.Context(), chatID, &req, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID чата"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

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

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	messageIDStr := c.Query("message_id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID сообщения"})
		return
	}

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
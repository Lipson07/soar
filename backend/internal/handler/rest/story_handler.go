package rest

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoryHandler struct {
	storyService service.StoryService
	userService  service.UserService
}

func NewStoryHandler(storyService service.StoryService, userService service.UserService) *StoryHandler {
	return &StoryHandler{
		storyService: storyService,
		userService:  userService,
	}
}

func (h *StoryHandler) getUserID(c *gin.Context) uuid.UUID {
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

func (h *StoryHandler) UploadStory(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Получаем информацию о пользователе из БД
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	userName := user.Username
	if userName == "" {
		userName = "Пользователь"
	}

	var avatarPtr *string
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		avatarPtr = user.AvatarURL
	}

	story, err := h.storyService.UploadStory(userID.String(), userName, avatarPtr, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, story)
}

func (h *StoryHandler) GetStories(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	stories, err := h.storyService.GetStories(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Если у сторис пустое имя или аватар, пробуем заполнить из БД
	for i, story := range stories {
		if story.UserName == "" || story.UserAvatar == nil {
			uid, err := uuid.Parse(story.UserID)
			if err == nil {
				u, err := h.userService.GetByID(c.Request.Context(), uid)
				if err == nil {
					if story.UserName == "" {
						stories[i].UserName = u.Username
						if stories[i].UserName == "" {
							stories[i].UserName = "Пользователь"
						}
					}
					if story.UserAvatar == nil && u.AvatarURL != nil && *u.AvatarURL != "" {
						stories[i].UserAvatar = u.AvatarURL
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, stories)
}

func (h *StoryHandler) MarkAsViewed(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	storyID := c.Param("id")
	if storyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story id is required"})
		return
	}

	if err := h.storyService.MarkStoryAsViewed(storyID, userID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

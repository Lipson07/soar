package rest

import (
	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(public, protected *gin.RouterGroup)
}

func SetupRouter(registrars ...RouteRegistrar) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	public := router.Group("/api")
	protected := router.Group("/api")
	protected.Use(AuthMiddleware())

	for _, registrar := range registrars {
		registrar.RegisterRoutes(public, protected)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
func (h *UserHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	public.POST("/register", h.Register)
	public.POST("/login", h.Login)

	users := protected.Group("/users")
	{
		users.GET("/profile", h.GetProfile)
		users.GET("/", h.GetAllUsers)
		users.GET("/search", h.SearchUsers)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.PUT("/profile/status", h.UpdateStatus)
	}
}

func (h *ChatHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	chats := protected.Group("/chats")
	{
		chats.POST("/private", h.CreatePrivateChat)
		chats.POST("/group", h.CreateGroupChat)
		chats.GET("/", h.GetUserChats)
		chats.GET("/all", h.GetAllChats)
		chats.GET("/:id", h.GetChatByID)
		chats.PUT("/:id", h.UpdateChat)
		chats.DELETE("/:id", h.DeleteChat)
	}
}

func (h *ParticipantHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	participants := protected.Group("/participants")
	{
		participants.POST("", h.AddParticipants)
		participants.DELETE("", h.RemoveParticipant)
		participants.POST("/leave", h.LeaveChat)
		participants.GET("", h.GetChatParticipants)
		participants.PUT("/role", h.UpdateRole)
		participants.PUT("/read", h.UpdateLastRead)
		participants.GET("/unread", h.GetUnreadCount)
	}
}

func (h *MessageHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	messages := protected.Group("/messages")
	{
		messages.POST("", h.SendMessage)
		messages.GET("", h.GetMessages)
		messages.PUT("", h.EditMessage)
		messages.DELETE("", h.DeleteMessage)
	}
}

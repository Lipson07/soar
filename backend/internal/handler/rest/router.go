package rest

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(public, protected *gin.RouterGroup)
}

func SetupRouter(registrars ...RouteRegistrar) *gin.Engine {
	router := gin.Default()

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	workDir, _ := os.Getwd()
	uploadsPath := filepath.Join(workDir, "uploads")
	imagesPath := filepath.Join(uploadsPath, "images")
	filesPath := filepath.Join(uploadsPath, "files")

	os.MkdirAll(imagesPath, 0755)
	os.MkdirAll(filesPath, 0755)

	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/uploads") {
			c.Header("Cache-Control", "public, max-age=31536000")
		}
		c.Next()
	})

	router.Static("/uploads", uploadsPath)

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
		users.GET("", h.GetAllUsers)
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
		chats.GET("", h.GetUserChats)
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
		messages.POST("/upload", h.UploadFile)
		messages.GET("", h.GetMessages)
		messages.PUT("", h.EditMessage)
		messages.DELETE("", h.DeleteMessage)
	}
}
func (h *SecurityHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	security := protected.Group("/security")
	{
		security.GET("/settings", h.GetSettings)
		security.PUT("/settings", h.UpdateSettings)
		security.POST("/2fa/setup", h.SetupTwoFactor)
		security.POST("/2fa/verify", h.VerifyTwoFactor)
		security.DELETE("/2fa", h.DisableTwoFactor)
		security.GET("/sessions", h.GetSessions)
		security.DELETE("/sessions/:sessionId", h.TerminateSession)
		security.DELETE("/sessions", h.TerminateAllOtherSessions)
		security.GET("/report", h.GetSecurityReport)
	}
}
func (h *FilesHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	files := protected.Group("/files")
	{
		files.GET("", h.GetFiles)
		files.DELETE("/*filepath", h.DeleteFile)
	}
}
func (h *CallHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	calls := protected.Group("/calls")
	{
		calls.POST("", h.StartCall)
		calls.GET("", h.GetUserCalls)
		calls.GET("/active", h.GetActiveCall)
		calls.POST("/:callId/accept", h.AcceptCall)
		calls.POST("/:callId/reject", h.RejectCall)
		calls.POST("/:callId/end", h.EndCall)
	}
}

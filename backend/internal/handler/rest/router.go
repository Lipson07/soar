package rest

import (
	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(public, protected *gin.RouterGroup)
}

func SetupRouter(registrars ...RouteRegistrar) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

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
		users.GET("/", h.GetAllUsers)
		users.GET("/:id", h.GetUser)
		users.GET("/search", h.SearchUsers)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

func (h *ChatHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	chats := protected.Group("/chats")
	{
		chats.POST("/", h.Create)
		chats.GET("/", h.GetAll)
		chats.GET("/:chat_id", h.GetByID)
		chats.GET("/by-name", h.GetByName)
		chats.PUT("/:chat_id", h.Update)
		chats.DELETE("/:chat_id", h.Delete)
	}
}
func (h *ChatMemberHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	members := protected.Group("/chats/:chat_id/members")
	{
		members.POST("", h.AddMember)
		members.POST("/bulk", h.AddMembers)
		members.GET("", h.GetChatMembers)
		members.GET("/count", h.GetMemberCount)
		members.GET("/:user_id", h.GetMember)
		members.PUT("/:user_id/role", h.UpdateMemberRole)
		members.DELETE("/:user_id", h.RemoveMember)
		members.POST("/:user_id/kick", h.KickMember)
	}

	protected.POST("/chats/:chat_id/read", h.UpdateLastRead)
	protected.POST("/chats/:chat_id/leave", h.LeaveChat)
	protected.GET("/users/chats", h.GetUserChats)
}

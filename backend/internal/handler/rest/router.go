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
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}



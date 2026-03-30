package app

import (
	"myapp/internal/handler/rest"
	"myapp/internal/repository/postgres"
	"myapp/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initDependencies(db *sqlx.DB) *gin.Engine {

	userRepo := postgres.NewUserRepository(db)
	chatRepo:=postgres.NewChatRepostory(db)

	userService := service.NewUserService(userRepo)
	chatService:=service.NewChatService(chatRepo)
	userHandler := rest.NewUserHandler(userService)
	chatHandler:=rest.NewChatHandler(chatService)
	router := rest.SetupRouter(
		userHandler,
		chatHandler,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	return router
}
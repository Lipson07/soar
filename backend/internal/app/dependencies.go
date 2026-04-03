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
	chatRepo := postgres.NewChatRepository(db)
	chatMemberRepo := postgres.NewChatMemberRepository(db)
	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(chatRepo, chatMemberRepo)
	chatMemberService := service.NewChatMemberService(chatMemberRepo, chatRepo, userRepo, db)
	userHandler := rest.NewUserHandler(userService)
	chatHandler := rest.NewChatHandler(chatService)
	chatMemberHandler := rest.NewChatMemberHandler(chatMemberService)
	router := rest.SetupRouter(
		userHandler,
		chatHandler,
		chatMemberHandler,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	return router
}

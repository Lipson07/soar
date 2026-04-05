package app

import (
	"backend/internal/handler/rest"
	"backend/internal/repository/postgres"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initDependencies(db *sqlx.DB) *gin.Engine {

userRepo := postgres.NewUserRepository(db)
	chatRepo := postgres.NewChatRepository(db)
	participantRepo := postgres.NewParticipantRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(chatRepo, participantRepo, userRepo, messageRepo)
	participantService := service.NewParticipantService(participantRepo, chatRepo, userRepo,messageRepo)
	messageService := service.NewMessageService(messageRepo, participantRepo, chatRepo)

	userHandler := rest.NewUserHandler(userService)
	chatHandler := rest.NewChatHandler(chatService)
	participantHandler := rest.NewParticipantHandler(participantService)
	messageHandler := rest.NewMessageHandler(messageService)

	router := rest.SetupRouter(
		userHandler,
		chatHandler,
		participantHandler,
		messageHandler,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	return router
}

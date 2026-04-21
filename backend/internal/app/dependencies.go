package app

import (
	"backend/internal/handler/rest"
	"backend/internal/handler/websocket"
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
	securityRepo := postgres.NewSecurityRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	callRepo := postgres.NewCallRepository(db)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(chatRepo, participantRepo, userRepo, messageRepo)
	participantService := service.NewParticipantService(participantRepo, chatRepo, userRepo, messageRepo)
	messageService := service.NewMessageService(messageRepo, participantRepo, chatRepo)
	securityService := service.NewSecurityService(securityRepo, sessionRepo)
	callService := service.NewCallService(callRepo, participantRepo, chatRepo)

	hub := websocket.NewHub()
	go hub.Run()

	wsHandler := websocket.NewWebSocketHandler(hub, callService, userService)

	userHandler := rest.NewUserHandler(userService)
	chatHandler := rest.NewChatHandler(chatService)
	participantHandler := rest.NewParticipantHandler(participantService)
	messageHandler := rest.NewMessageHandler(messageService)
	securityHandler := rest.NewSecurityHandler(securityService)
	filesHandler := rest.NewFilesHandler()
	callHandler := rest.NewCallHandler(callService)

	router := rest.SetupRouter(
		userHandler,
		chatHandler,
		participantHandler,
		messageHandler,
		securityHandler,
		filesHandler,
		callHandler,
	)

	router.GET("/ws", wsHandler.HandleWebSocket)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	return router
}

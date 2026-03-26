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


	userService := service.NewUserService(userRepo)
	
	userHandler := rest.NewUserHandler(userService)

	router := rest.SetupRouter(
		userHandler,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	return router
}
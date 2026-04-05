package main

import (
	_ "backend/docs"
	"backend/internal/app"
	"log"
)

// @title           MyApp API
// @version         1.0.0
// @description     API для управления пользователями и проектами
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      backend-1-8qiq.onrender.com
// @BasePath  /api
// @schemes   https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal("Ошибка инициализации приложения:", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("Ошибка запуска приложения:", err)
	}
}
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/repository/postgres"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

	type App struct {
		cfg    *config.Config
		db     *sqlx.DB
		router *gin.Engine
		server *http.Server
	}

	func New() (*App, error) {

		cfg := config.Load()
		log.Println("Конфигурация загружена")

		db, err := postgres.NewDB(cfg)
		if err != nil {
			return nil, err
		}
		log.Println("Подключено к PostgreSQL")

		router := initDependencies(db)

		server := &http.Server{
			Addr:         ":" + cfg.HTTP_PORT,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		return &App{
			cfg:    cfg,
			db:     db,
			router: router,
			server: server,
		}, nil
	}

	func (a *App) Run() error {

		go func() {
			log.Printf("Сервер запущен на порту %s", a.cfg.HTTP_PORT)
			log.Printf("Swagger UI доступен по адресу: http://localhost:%s/swagger/index.html", a.cfg.HTTP_PORT)
			if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Ошибка сервера: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Останавливаем сервер...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			return err
		}

		if err := a.db.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}

		log.Println("Сервер успешно остановлен")
		return nil
	}
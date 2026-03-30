package rest

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	RegisterRoutes(public, protected *gin.RouterGroup)
}

type ChatHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetByName(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RegisterRoutes(public, protected *gin.RouterGroup)
} 	
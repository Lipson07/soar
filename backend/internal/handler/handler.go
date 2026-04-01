package rest

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUser(c *gin.Context)
	SearchUsers(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type ChatHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetByName(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type ChatMemberHandler interface {
	AddMember(c *gin.Context)
	AddMembers(c *gin.Context)
	GetChatMembers(c *gin.Context)
	GetMember(c *gin.Context)
	UpdateMemberRole(c *gin.Context)
	UpdateLastRead(c *gin.Context)
	RemoveMember(c *gin.Context)
	LeaveChat(c *gin.Context)
	KickMember(c *gin.Context)
	GetUserChats(c *gin.Context)
	GetMemberCount(c *gin.Context)
}
